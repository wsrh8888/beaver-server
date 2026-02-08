package logic

import (
	"context"
	"time"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
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
	// 1. 获取邀请者（发起人）信息
	callerInfo, _ := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: req.UserID})
	var callerUserInfo map[string]string
	if callerInfo != nil && callerInfo.GetUserInfo() != nil {
		callerUserInfo = map[string]string{
			"userId":   req.UserID,
			"nickName": callerInfo.GetUserInfo().NickName,
			"avatar":   callerInfo.GetUserInfo().Avatar,
		}
	}

	// 2. 依次登记状态并发送信令
	for _, targetID := range req.TargetIds {
		// 登记为待接听状态
		_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
			RoomId: req.RoomID,
			UserId: targetID,
			Status: 1, // 1-待接听
		})

		// 批量通过 WebSocket 发送强通知 RTC 信令
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd,
			wsCommandConst.CALL,
			wsTypeConst.CallReceive,
			req.UserID,
			targetID,
			map[string]interface{}{
				"type":           "RTC_INVITE",
				"roomId":         req.RoomID,
				"callerId":       req.UserID,
				"callType":       2, // 邀请必然发生在群聊上下文中
				"callerUserInfo": callerUserInfo,
				"timestamp":      time.Now().Unix(),
			},
			"group_call", // 群邀请信令会话标记
		)
	}

	return &types.InviteCallMemberRes{}, nil
}
