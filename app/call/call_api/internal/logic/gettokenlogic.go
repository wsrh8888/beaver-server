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

	"github.com/livekit/protocol/auth"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 接听通话并获取令牌
func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *types.GetCallTokenReq) (resp *types.GetCallTokenRes, err error) {
	// 1. 校验会话和参与者身份
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		return nil, errors.New("通话会话不存在")
	}

	isParticipant := false
	for _, pid := range session.ParticipantIds {
		if pid == req.UserID {
			isParticipant = true
			break
		}
	}
	if !isParticipant {
		return nil, errors.New("无权进入该通话")
	}

	// 2. 更新状态为已接听
	_, err = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
		RoomId: req.RoomID,
		UserId: req.UserID,
		Status: 2, // 2-已接听
	})
	if err != nil {
		return nil, err
	}

	// 3. 生成令牌
	token, err := l.generateToken(req.UserID, req.RoomID)
	if err != nil {
		return nil, err
	}

	// 4. 发送信令告知发起者有人接听
	go l.sendAcceptSignal(req.UserID, session.CallerId, req.RoomID)

	return &types.GetCallTokenRes{
		RoomToken:  token,
		LiveKitUrl: l.svcCtx.Config.LiveKit.Host,
	}, nil
}

func (l *GetTokenLogic) generateToken(userID, roomID string) (string, error) {
	at := auth.NewAccessToken(l.svcCtx.Config.LiveKit.ApiKey, l.svcCtx.Config.LiveKit.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomID,
	}
	at.AddGrant(grant).SetIdentity(userID).SetValidFor(time.Hour)
	return at.ToJWT()
}

func (l *GetTokenLogic) sendAcceptSignal(acceptorID, callerID, roomID string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type":   "RTC_ACCEPTED",
		"user":   acceptorID,
		"roomId": roomID,
	})

	_, err := l.svcCtx.ChatRpc.SendMsg(context.Background(), &chat_rpc.SendMsgReq{
		UserId:         acceptorID,
		ConversationId: l.getConversationID(acceptorID, callerID),
		Msg: &chat_rpc.Msg{
			Type: 7, // 7:通知消息/信令
			NotificationMsg: &chat_rpc.NotificationMsg{
				Type:   101, // RTC_ACCEPTED
				Actors: []string{acceptorID},
			},
			TextMsg: &chat_rpc.TextMsg{
				Content: string(payload),
			},
		},
	})
	if err != nil {
		logx.Errorf("发送 RTC_ACCEPTED 信令失败: %v", err)
	}
}

func (l *GetTokenLogic) getConversationID(u1, u2 string) string {
	if u1 < u2 {
		return u1 + ":" + u2
	}
	return u2 + ":" + u1
}
