package logic

import (
	"context"
	"strconv"

	"beaver/app/chat/chat_models"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	"beaver/core/coreonline"
	"beaver/core/corepush"

	"github.com/zeromicro/go-zero/core/logx"
)

func (l *SendMsgLogic) sendOfflinePushIfNeeded(
	conversationID, senderID string,
	chatModel chat_models.ChatMessage,
	recipientIDs []string,
) {
	if l.svcCtx.PushSender == nil || !l.svcCtx.PushSender.Enabled() {
		return
	}

	sender, err := l.getSenderInfo(chatModel)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("离线推送获取发送者失败: %v", err)
		return
	}

	title := sender.NickName
	if title == "" {
		title = "新消息"
	}
	body := chatModel.MsgPreview
	if body == "" {
		body = "你收到一条新消息"
	}

	data := map[string]string{
		"type":           "chat_message",
		"conversationId": conversationID,
		"messageId":      chatModel.MessageID,
		"seq":            strconv.FormatInt(chatModel.Seq, 10),
		"senderId":       senderID,
	}

	for _, recipientID := range recipientIDs {
		if recipientID == senderID {
			continue
		}
		if coreonline.IsOnline(l.svcCtx.Redis, recipientID) {
			continue
		}

		var userConvo chat_models.ChatUserConversation
		if err := l.svcCtx.DB.Where("user_id = ? AND conversation_id = ?", recipientID, conversationID).
			First(&userConvo).Error; err == nil && userConvo.IsMuted {
			continue
		}

		res, err := l.svcCtx.NotificationRpc.ListPushTokens(context.Background(), &notification_rpc.ListPushTokensReq{
			UserId: recipientID,
		})
		if err != nil || len(res.Tokens) == 0 {
			continue
		}

		tokens := make([]corepush.PushToken, 0, len(res.Tokens))
		for _, t := range res.Tokens {
			tokens = append(tokens, corepush.PushToken{
				DeviceID:     t.DeviceId,
				PushToken:    t.PushToken,
				PushPlatform: t.PushPlatform,
			})
		}

		l.svcCtx.PushSender.SendToTokens(context.Background(), tokens, corepush.Message{
			Title: title,
			Body:  body,
			Data:  data,
		})
	}
}
