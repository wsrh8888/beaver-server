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

	// 1. 根据通话类型处理目标
	if req.CallType == 1 || req.CallType == 2 {
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
	} else if req.CallType == 3 || req.CallType == 4 {
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

	// 4. 生成 LiveKit Token
	token, err := l.generateToken(req.UserID, roomID)
	if err != nil {
		return nil, err
	}

	// 5. 发送信令 (批量发送)
	go l.sendInviteSignals(req.UserID, req.TargetId, roomID, req.CallType, targetIDs)

	return &types.StartCallRes{
		RoomID:     roomID,
		RoomToken:  token,
		LiveKitUrl: l.svcCtx.Config.LiveKit.Host,
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

func (l *StartCallLogic) sendInviteSignals(callerID, conversationID, roomID string, callType int8, targetIDs []string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type":     "RTC_INVITE",
		"roomId":   roomID,
		"callerId": callerID,
		"callType": callType,
	})

	// 1. 通过 ChatRpc 发送一条信令消息到会话（群聊或私聊）
	_, err := l.svcCtx.ChatRpc.SendMsg(context.Background(), &chat_rpc.SendMsgReq{
		UserId:         callerID,
		ConversationId: l.getConversationID(callerID, conversationID, callType),
		Msg: &chat_rpc.Msg{
			Type: 7, // 7:通知消息/信令
			NotificationMsg: &chat_rpc.NotificationMsg{
				Type:   100, // 自定义 RTC 类型
				Actors: []string{callerID},
			},
			TextMsg: &chat_rpc.TextMsg{
				Content: string(payload),
			},
		},
	})
	if err != nil {
		logx.Errorf("发送 RTC_INVITE 聊天信令失败: %v", err)
	}

	// 2. 批量通过 WebSocket 发送强通知 RTC 信令
	for _, targetID := range targetIDs {
		ajax.SendMessageToWs(l.svcCtx.Config.Etcd,
			wsCommandConst.CALL,
			wsTypeConst.CallReceive,
			callerID,
			targetID,
			map[string]interface{}{
				"type":      "RTC_INVITE",
				"roomId":    roomID,
				"callerId":  callerID,
				"callType":  callType,
				"timestamp": time.Now().Unix(),
			},
			l.getConversationID(callerID, conversationID, callType),
		)
	}
}

func (l *StartCallLogic) getConversationID(callerID, targetID string, callType int8) string {
	if callType == 3 || callType == 4 {
		return targetID // 群聊会话ID就是群ID
	}
	// 私聊会话ID拼装
	if callerID < targetID {
		return callerID + ":" + targetID
	}
	return targetID + ":" + callerID
}
