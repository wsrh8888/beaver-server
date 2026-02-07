package logic

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type HangupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 主动挂断或拒绝通话
func NewHangupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HangupLogic {
	return &HangupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HangupLogic) Hangup(req *types.HangupCallReq) (resp *types.HangupCallRes, err error) {
	// 1. 获取会话信息，找到对方
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: req.RoomID})
	if err != nil {
		return nil, errors.New("通话会话不存在")
	}

	// 2. 更新 RPC 状态为已结束
	_, err = l.svcCtx.CallRpc.FinalizeSession(l.ctx, &call_rpc.FinalizeSessionReq{
		RoomId: req.RoomID,
		Status: 3, // 3-已结束
	})
	if err != nil {
		return nil, err
	}

	// 3. 更新参与者状态
	_, err = l.svcCtx.CallRpc.UpdateParticipantStatus(l.ctx, &call_rpc.UpdateParticipantStatusReq{
		RoomId: req.RoomID,
		UserId: req.UserID,
		Status: 5, // 5-已挂断
	})

	// 4. 发送信令告知对方挂断
	for _, pid := range session.ParticipantIds {
		if pid != req.UserID {
			go l.sendHangupSignal(req.UserID, pid, req.RoomID)
		}
	}

	return &types.HangupCallRes{Success: true}, nil
}

func (l *HangupLogic) sendHangupSignal(hanguperID, targetID, roomID string) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type":   "RTC_HANGUP",
		"user":   hanguperID,
		"roomId": roomID,
	})

	_, err := l.svcCtx.ChatRpc.SendMsg(context.Background(), &chat_rpc.SendMsgReq{
		UserId:         hanguperID,
		ConversationId: l.getConversationID(hanguperID, targetID),
		Msg: &chat_rpc.Msg{
			Type: 7, // 7:通知消息/信令
			NotificationMsg: &chat_rpc.NotificationMsg{
				Type:   102, // RTC_HANGUP
				Actors: []string{hanguperID},
			},
			TextMsg: &chat_rpc.TextMsg{
				Content: string(payload),
			},
		},
	})
	if err != nil {
		logx.Errorf("发送 RTC_HANGUP 信令失败: %v", err)
	}

	// 2. 直接通过 WebSocket 发送 RTC 信令
	ajax.SendMessageToWs(l.svcCtx.Config.Etcd,
		wsCommandConst.CALL,
		wsTypeConst.CallReceive,
		hanguperID,
		targetID,
		map[string]interface{}{
			"type":   "RTC_HANGUP",
			"user":   hanguperID,
			"roomId": roomID,
		},
		l.getConversationID(hanguperID, targetID),
	)
}

func (l *HangupLogic) getConversationID(u1, u2 string) string {
	if u1 < u2 {
		return u1 + ":" + u2
	}
	return u2 + ":" + u1
}
