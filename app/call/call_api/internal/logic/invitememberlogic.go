package logic

import (
	"context"
	"time"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群聊中邀请成员入场
func NewInviteMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteMemberLogic {
	return &InviteMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteMemberLogic) InviteMember(req *types.InviteCallMemberReq) (resp *types.InviteCallMemberRes, err error) {
	// 1. 获取会话信息 (主要为了拿到 ConversationID 和 CallType)
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		return nil, err
	}

	// 2. 获取邀请者（发起人）信息
	callerInfo, _ := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: req.UserID})
	var callerUserInfo map[string]string
	if callerInfo != nil && callerInfo.GetUserInfo() != nil {
		callerUserInfo = map[string]string{
			"userId":   req.UserID,
			"nickName": callerInfo.GetUserInfo().NickName,
			"avatar":   callerInfo.GetUserInfo().Avatar,
		}
	}

	// 3. 依次登记状态并发送信令
	for _, targetID := range req.TargetIds {
		// 登记为待接听状态
		_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
			RoomId: req.RoomID,
			UserId: targetID,
			Status: 1, // 1-待接听 (ParticipantStatusCalling)
		})

		// 通过 WebSocket 发送 RTC_INVITE 信令
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd,
			wsCommandConst.CALL,
			wsTypeConst.CallReceive,
			req.UserID,
			targetID,
			map[string]interface{}{
				"type":           call_models.SignalInvite,
				"roomId":         req.RoomID,
				"callerId":       req.UserID,
				"callType":       session.CallType,
				"callerUserInfo": callerUserInfo,
				"timestamp":      time.Now().Unix(),
			},
			session.ConversationId,
		)

		// 4. 开启超时处理定时器 (60秒未接听则自动设为超时)
		l.startTimeoutTimer(req.RoomID, targetID)
	}

	return &types.InviteCallMemberRes{}, nil
}

// startTimeoutTimer 异步计时器：如果用户在规定时间内未接听，则自动更新状态为超时
func (l *InviteMemberLogic) startTimeoutTimer(roomID, userID string) {
	time.AfterFunc(60*time.Second, func() {
		// 使用 context.Background()，因为原始请求的上下文会因接口返回而取消
		ctx := context.Background()

		// 1. 确认用户当前状态
		participants, err := l.svcCtx.CallRpc.GetParticipants(ctx, &call_rpc.GetParticipantsReq{RoomId: roomID})
		if err != nil {
			return
		}

		isStillCalling := false
		for _, p := range participants.Participants {
			if p.UserId == userID && p.Status == 1 { // 1 代表 Calling
				isStillCalling = true
				break
			}
		}

		// 2. 如果依然是待接听状态，则变更为超时
		if isStillCalling {
			_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(ctx, &call_rpc.UpdateParticipantStatusReq{
				RoomId: roomID,
				UserId: userID,
				Status: 4, // 4 代表 Timeout
			})
		}
	})
}
