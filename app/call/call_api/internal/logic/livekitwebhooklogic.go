package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/livekit/protocol/livekit"
	"github.com/zeromicro/go-zero/core/logx"
)

type LiveKitWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// LiveKit 服务器回调 (需在网关配置白名单)
func NewLiveKitWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveKitWebhookLogic {
	return &LiveKitWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveKitWebhookLogic) LiveKitWebhook(req *types.LiveKitWebhookReq) (resp *types.LiveKitWebhookRes, err error) {
	// 1. 校验来源 (TODO: 使用 LiveKit SDK 校验)
	// 在生产环境中应启用签名校验

	// 2. 解析事件
	var event livekit.WebhookEvent
	if err := json.Unmarshal(req.Body, &event); err != nil {
		l.Errorf("解析 Webhook 事件失败: %v", err)
		return nil, err
	}

	roomID := event.Room.Name
	l.Infof("收到 LiveKit Webhook 事件: %s, Room: %s", event.Event, roomID)

	switch event.Event {
	case "room_finished":
		// 通话彻底结束，数据沉淀
		var duration int32
		if event.Room != nil {
			duration = int32(event.CreatedAt - event.Room.CreationTime)
		}

		_, err = l.svcCtx.CallRpc.FinalizeSession(l.ctx, &call_rpc.FinalizeSessionReq{
			RoomId:   roomID,
			Duration: duration,
			Status:   3, // 3-已结束
		})
		if err != nil {
			l.Errorf("FinalizeSession 失败: %v", err)
		}

		// 写入聊天历史
		l.sendEndMessage(roomID, duration)
	}

	return &types.LiveKitWebhookRes{}, nil
}

func (l *LiveKitWebhookLogic) sendEndMessage(roomID string, duration int32) {
	session, err := l.svcCtx.CallRpc.GetSession(l.ctx, &call_rpc.GetSessionReq{RoomId: roomID})
	if err != nil {
		return
	}

	mins := duration / 60
	secs := duration % 60
	content := fmt.Sprintf("[通话结束] 时长 %02d:%02d", mins, secs)

	// 给参与者发消息
	for _, pid := range session.ParticipantIds {
		if pid == session.CallerId {
			continue
		}
		_, _ = l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
			UserId:         session.CallerId,
			ConversationId: l.getConversationID(session.CallerId, pid),
			Msg: &chat_rpc.Msg{
				Type: 1, // 1:文本
				TextMsg: &chat_rpc.TextMsg{
					Content: content,
				},
			},
		})
	}
}

func (l *LiveKitWebhookLogic) getConversationID(u1, u2 string) string {
	if u1 < u2 {
		return u1 + ":" + u2
	}
	return u2 + ":" + u1
}
