package logic

import (
	"context"
	"time"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	user_rpc "beaver/app/user/user_rpc/types/user_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/conversation"

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
	// 1. 生成房间ID (同时也作为通话记录在聊天中的初始消息ID)
	roomID := uuid.New().String()

	// 2. 解析会话ID获取目标ID (使用工具类统一处理)
	convID := req.ConversationId
	targetID := conversation.GetTargetIDByConversation(convID, req.UserID)

	// 3. 调用 RPC 创建会话
	_, err = l.svcCtx.CallRpc.CreateSession(l.ctx, &call_rpc.CreateSessionReq{
		RoomId:         roomID,
		CallerId:       req.UserID,
		TargetId:       targetID,
		CallType:       int32(req.CallType),
		MessageId:      roomID, // 使用 roomID 作为锚点消息ID
		ConversationId: convID,
	})
	if err != nil {
		return nil, err
	}

	// 4. 生成令牌
	token, err := l.generateToken(req.UserID, roomID)
	if err != nil {
		return nil, err
	}

	// 5. 同步发送第一条“进行中”聊天消息 (作为存证和后期状态补丁的锚点)
	l.sendInitialCallMessage(req.UserID, convID, roomID, roomID, req.CallType)

	// 6. 异步发送 WebSocket 实时弹窗信号 (仅私聊受邀方需要唤起界面；群聊由用户按需进入)
	if req.CallType == 1 {
		go l.sendInviteSignals(req.UserID, targetID, convID, roomID, req.CallType)
		// 开启超时处理定时器
		l.startTimeoutTimer(roomID, targetID)
	}

	return &types.StartCallRes{
		RoomID:     roomID,
		RoomToken:  token,
		LiveKitUrl: l.svcCtx.Config.LiveKit.Host,
		MessageID:  roomID,
	}, nil
}

// startTimeoutTimer 异步计时器：如果用户在规定时间内未接听，则自动更新状态为超时
func (l *StartCallLogic) startTimeoutTimer(roomID, userID string) {
	time.AfterFunc(30*time.Second, func() {
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

func (l *StartCallLogic) generateToken(userID, roomID string) (string, error) {
	at := auth.NewAccessToken(l.svcCtx.Config.LiveKit.ApiKey, l.svcCtx.Config.LiveKit.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}
	at.AddGrant(grant).SetIdentity(userID).SetValidFor(time.Hour)
	return at.ToJWT()
}

func (l *StartCallLogic) sendInitialCallMessage(callerID, convID, roomID, messageID string, callType int8) {
	// 1. 发送第一条“进行中”记录 (这就是我们的永久锚点)
	_, _ = l.svcCtx.ChatRpc.SendMsg(context.Background(), &chat_rpc.SendMsgReq{
		UserId:         callerID,
		ConversationId: convID,
		MessageId:      messageID,
		Msg: &chat_rpc.Msg{
			Type: 9, // 9: 音视频通话
			CallMsg: &chat_rpc.CallMsg{
				RoomId:   roomID,
				CallType: int32(callType),
				Status:   1, // 1: 进行中
			},
		},
	})
}

func (l *StartCallLogic) sendInviteSignals(callerID, targetID, convID, roomID string, callType int8) {
	// 1. 获取发起者基本信息用于被叫方显示
	var callerUserInfo map[string]string
	callerInfo, _ := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{UserID: callerID})
	if callerInfo != nil && callerInfo.GetUserInfo() != nil {
		callerUserInfo = map[string]string{
			"userId":   callerID,
			"nickName": callerInfo.GetUserInfo().NickName,
			"avatar":   callerInfo.GetUserInfo().Avatar,
		}
	}

	// 2. 通过 RocketMQ 异步发送 WebSocket 实时弹窗信号 (纯信令，不入库)
	payload := map[string]interface{}{
		"command":  wsCommandConst.CALL,
		"type":     wsTypeConst.CallReceive,
		"senderId": callerID,
		"targetId": targetID,
		"body": map[string]interface{}{
			"type":           call_models.SignalInvite,
			"roomId":         roomID,
			"callerId":       callerID,
			"callType":       callType,
			"callerUserInfo": callerUserInfo,
			"timestamp":      time.Now().Unix(),
		},
		"conversationId": convID,
	}
	l.svcCtx.RocketMQ.SendMessage(l.ctx, mqwsconst.MqTopicWs, payload)
}
