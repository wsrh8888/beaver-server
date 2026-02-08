package logic

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/google/uuid"
	"github.com/livekit/protocol/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type StartCallLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发起音视频通话
func NewStartCallLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartCallLogic {
	return &StartCallLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StartCallLogic) StartCall(req *types.StartCallReq) (resp *types.StartCallRes, err error) {
	var targetIDs []string

	// 1. 根据通话类型处理目标 (1-私聊, 2-群聊)
	if req.CallType == 1 {
		// 单聊
		_, err = l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: req.TargetId})
		if err != nil {
			return nil, errors.New("目标用户不存在")
		}

		// 检查对方是否忙碌
		statusRes, err := l.svcCtx.CallRpc.GetUserStatus(l.ctx, &call_rpc.GetUserStatusReq{UserId: req.TargetId})
		if err == nil && statusRes.IsBusy {
			return nil, errors.New("对方正在通话中")
		}
		targetIDs = append(targetIDs, req.TargetId)
	} else if req.CallType == 2 {
		// 群聊
		memberRes, err := l.svcCtx.GroupRpc.GetGroupMembers(l.ctx, &group_rpc.GetGroupMembersReq{GroupID: req.TargetId})
		if err != nil {
			return nil, errors.New("获取群成员失败")
		}
		for _, m := range memberRes.Members {
			if m.UserID != req.UserID {
				targetIDs = append(targetIDs, m.UserID)
			}
		}
	} else {
		return nil, errors.New("非法的通话类型")
	}

	// 2. 生成唯一 RoomID
	roomID := uuid.New().String()

	// 3. 调用 RPC 创建会话
	_, err = l.svcCtx.CallRpc.CreateSession(l.ctx, &call_rpc.CreateSessionReq{
		RoomId:   roomID,
		CallerId: req.UserID,
		TargetId: req.TargetId,
		CallType: int32(req.CallType),
	})
	if err != nil {
		return nil, err
	}

	// 3.5 参与者登记逻辑
	// 无论什么通话，发起者自己首先进入房间 (joined)
	_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
		RoomId: roomID,
		UserId: req.UserID,
		Status: 2, // 2-已接听 (joined)
	})

	// 如果是私聊 (Type 1)，需要将目标用户设为 calling 状态，以便对方振铃和本地显示“呼叫中”
	if req.CallType == 1 {
		_, _ = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
			RoomId: roomID,
			UserId: req.TargetId,
			Status: 1, // 1-待接听 (calling)
		})
	}
	// 如果是群聊 (Type 2)，不自动登记其他成员。其他人通过 AddMember 加入或由管理员邀请。

	// 4. 获取发起者用户信息（用于被叫方显示）
	callerInfo, _ := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: req.UserID})
	var callerUserInfo map[string]string
	if callerInfo != nil && callerInfo.GetUserInfo() != nil {
		callerUserInfo = map[string]string{
			"name":   callerInfo.GetUserInfo().NickName,
			"avatar": callerInfo.GetUserInfo().Avatar,
		}
	}

	// 5. 生成 LiveKit Token
	token, err := l.generateToken(req.UserID, roomID)
	if err != nil {
		return nil, err
	}

	// 通话邀请逻辑
	// 如果是私聊，发送邀请信令。如果是群聊，不做任何主动通知（用户要求纯净开场）。
	if req.CallType == 1 {
		go l.sendInviteSignals(req.UserID, req.TargetId, roomID, req.CallType, targetIDs, callerUserInfo)
	}

	// 7. 统一通过 RPC 获取权威名单快照
	rpcResp, err := l.svcCtx.CallRpc.GetParticipants(l.ctx, &call_rpc.GetParticipantsReq{
		RoomId: roomID,
	})

	participants := make([]types.Participant, 0)
	if err == nil {
		for _, p := range rpcResp.Participants {
			// 过滤：只返回当前在房间或正在呼叫的人 (1-待接听/calling, 2-已接听/joined)
			// 拒绝、离开、忙线的人在初始化快照中应该被排除
			if p.Status != 1 && p.Status != 2 {
				continue
			}

			status := "calling"
			if p.Status == 2 {
				status = "joined"
			}
			participants = append(participants, types.Participant{
				UserID: p.UserId,
				Status: status,
			})
		}
	}

	return &types.StartCallRes{
		RoomID:       roomID,
		RoomToken:    token,
		LiveKitUrl:   l.svcCtx.Config.LiveKit.Host,
		Participants: participants,
	}, nil
}

func (l *StartCallLogic) generateToken(userID, roomID string) (string, error) {
	at := auth.NewAccessToken(l.svcCtx.Config.LiveKit.ApiKey, l.svcCtx.Config.LiveKit.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}
	at.AddGrant(grant).SetIdentity(userID).SetValidFor(time.Hour)
	return at.ToJWT()
}

func (l *StartCallLogic) sendInviteSignals(callerID, conversationID, roomID string, callType int8, targetIDs []string, callerUserInfo map[string]string) {
	// 1. 通过 ChatRpc 发送信令消息 (私聊)
	payload, _ := json.Marshal(map[string]interface{}{
		"type":     "RTC_INVITE",
		"roomId":   roomID,
		"callType": callType,
		"event":    "INVITE",
	})

	_, _ = l.svcCtx.ChatRpc.SendMsg(context.Background(), &chat_rpc.SendMsgReq{
		UserId:         callerID,
		ConversationId: l.getConversationID(callerID, conversationID, callType),
		MessageId:      uuid.New().String(),
		Msg: &chat_rpc.Msg{
			Type: 7,
			NotificationMsg: &chat_rpc.NotificationMsg{
				Type:   103,
				Actors: []string{callerID},
			},
			TextMsg: &chat_rpc.TextMsg{
				Content: string(payload),
			},
		},
	})

	// 2. 通过 WebSocket 发送弹窗信令 (仅针对私聊 Target)
	for _, targetID := range targetIDs {
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd,
			wsCommandConst.CALL,
			wsTypeConst.CallReceive,
			callerID,
			targetID,
			map[string]interface{}{
				"type":           "RTC_INVITE",
				"roomId":         roomID,
				"callerId":       callerID,
				"callType":       callType,
				"callerUserInfo": callerUserInfo,
				"timestamp":      time.Now().Unix(),
			},
			l.getConversationID(callerID, targetID, callType),
		)
	}
}

func (l *StartCallLogic) getConversationID(callerID, targetID string, callType int8) string {
	if callType == 2 { // 2-群聊
		return targetID // 群聊会话ID就是群ID
	}
	// 私聊会话ID拼装 (与前端 private_uid1_uid2 保持一致方便解析)
	if callerID < targetID {
		return "private_" + callerID + "_" + targetID
	}
	return "private_" + targetID + "_" + callerID
}
