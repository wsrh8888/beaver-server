package logic

import (
	"context"
	"encoding/json"
	"strings"

	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/openevent"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/common/models/ctype"
	"beaver/utils/conversation"

	"github.com/zeromicro/go-zero/core/logx"
)

type robotWebhookPusher struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func newRobotWebhookPusher(ctx context.Context, svcCtx *svc.ServiceContext) *robotWebhookPusher {
	return &robotWebhookPusher{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (p *robotWebhookPusher) tryPush(in *chat_rpc.SendMsgReq, msg *ctype.Msg) {
	defer func() {
		if r := recover(); r != nil {
			p.Errorf("Robot Webhook 推送 panic: %v", r)
		}
	}()

	conversationType, userIDs := conversation.ParseConversationWithType(in.ConversationId)
	eventBody := p.buildEventBody(in, msg)

	switch conversationType {
	case 1:
		p.pushPrivateChat(in.UserId, userIDs, eventBody)
	case 2:
		p.pushGroupAt(in.UserId, in.ConversationId, msg, eventBody)
	}
}

func (p *robotWebhookPusher) pushPrivateChat(senderID string, userIDs []string, event map[string]interface{}) {
	if len(userIDs) != 2 {
		return
	}
	var peerID string
	for _, uid := range userIDs {
		if uid != senderID {
			peerID = uid
			break
		}
	}
	if peerID == "" {
		return
	}

	res, err := p.svcCtx.OpenRpc.GetRobotByUserID(p.ctx, &open_rpc.GetRobotByUserIDReq{RobotUserId: peerID})
	if err != nil || res == nil || !res.Found || !res.EnableSingleChat {
		return
	}

	p.dispatch(res.AppId, openevent.EventIMMessageReceive, event)
}

func (p *robotWebhookPusher) pushGroupAt(senderID, conversationID string, msg *ctype.Msg, event map[string]interface{}) {
	if msg == nil || len(msg.AtUserIDs) == 0 {
		return
	}

	for _, atUserID := range msg.AtUserIDs {
		res, err := p.svcCtx.OpenRpc.GetRobotByUserID(p.ctx, &open_rpc.GetRobotByUserIDReq{RobotUserId: atUserID})
		if err != nil || res == nil || !res.Found || !res.EnableAtMention {
			continue
		}
		groupEvent := copyEventMap(event)
		groupEvent["group_id"] = conversation.GetTargetIDByConversation(conversationID, senderID)
		p.dispatch(res.AppId, openevent.EventIMMessageReceiveGroup, groupEvent)
	}
}

func (p *robotWebhookPusher) dispatch(appID, eventType string, event map[string]interface{}) {
	body, err := json.Marshal(event)
	if err != nil {
		p.Errorf("Robot Webhook 事件序列化失败: %v", err)
		return
	}
	_, err = p.svcCtx.OpenRpc.DispatchPlatformEvent(p.ctx, &open_rpc.DispatchPlatformEventReq{
		AppId:     appID,
		EventType: eventType,
		EventJson: string(body),
	})
	if err != nil {
		p.Errorf("Robot Webhook 推送 RPC 失败: app=%s event=%s err=%v", appID, eventType, err)
	}
}

func (p *robotWebhookPusher) buildEventBody(in *chat_rpc.SendMsgReq, msg *ctype.Msg) map[string]interface{} {
	event := map[string]interface{}{
		"sender_id":       in.UserId,
		"conversation_id": in.ConversationId,
		"message_id":      in.MessageId,
		"msg_type":        msgTypeName(msg),
		"content":         extractTextContent(msg),
		"mentions":        msg.AtUserIDs,
	}
	if event["mentions"] == nil {
		event["mentions"] = []string{}
	}
	return event
}

func msgTypeName(msg *ctype.Msg) string {
	if msg == nil {
		return "text"
	}
	switch msg.Type {
	case ctype.TextMsgType:
		return "text"
	case ctype.MarkdownMsgType:
		return "markdown"
	case ctype.ImageMsgType:
		return "image"
	default:
		return "unknown"
	}
}

func extractTextContent(msg *ctype.Msg) string {
	if msg == nil {
		return ""
	}
	if msg.TextMsg != nil {
		return msg.TextMsg.Content
	}
	if msg.MarkdownMsg != nil {
		return msg.MarkdownMsg.Content
	}
	return ""
}

func copyEventMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func isRobotPeer(ctx context.Context, openRpc open_rpc.OpenClient, peerUserID string) bool {
	if peerUserID == "" {
		return false
	}
	res, err := openRpc.GetRobotByUserID(ctx, &open_rpc.GetRobotByUserIDReq{RobotUserId: peerUserID})
	return err == nil && res != nil && res.Found
}

func isOpenRobotSender(deviceID string) bool {
	return strings.EqualFold(deviceID, "open_robot")
}
