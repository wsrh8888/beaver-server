package chat_message

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"beaver/app/open/open_models"
	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// handleBotStreaming 处理 Bot 流式回复（大厂标准）
func handleBotStreaming(
	ctx context.Context,
	db *gorm.DB,
	client *ws_conn.Client,
	conversationID string,
	userMessage string,
	senderID string,
) {
	// 1. 检查会话是否关联 Bot
	var botApp open_models.OpenApp
	err := db.Where("bot_conversation_id = ?", conversationID).First(&botApp).Error
	if err != nil {
		// 不是 Bot 对话，直接返回
		return
	}

	// 2. 获取 Webhook 配置
	var webhookConfig open_models.OpenWebhookConfig
	err = db.Where("app_id = ? AND event_type = ? AND status = ?",
		botApp.AppID, "message.receive", 1).First(&webhookConfig).Error
	if err != nil {
		fmt.Printf("未找到 Webhook 配置: app_id=%s\n", botApp.AppID)
		return
	}

	// 3. 构建 Webhook payload
	payload := map[string]interface{}{
		"event":     "message.receive",
		"timestamp": fmt.Sprintf("%d", getCurrentTimestamp()),
		"app_id":    botApp.AppID,
		"message": map[string]interface{}{
			"id":              generateMessageID(),
			"conversation_id": conversationID,
			"sender_id":       senderID,
			"content":         userMessage,
			"msg_type":        "text",
		},
		"bot": map[string]interface{}{
			"id":   botApp.BotUserID,
			"name": botApp.Name,
		},
	}

	// 4. 准备消息 ID 和累积器
	botMsgID := generateMessageID()
	var fullContent strings.Builder
	var messageFormat string = "text" // 默认纯文本

	// 5. 创建流式处理器
	handler := &open_models.BotStreamHandler{
		WebhookURL: webhookConfig.TargetURL,
		Secret:     webhookConfig.Secret,
		Timeout:    webhookConfig.Timeout,
		Payload:    payload,
		OnChunk: func(chunk *open_models.BotChunk) error {
			// 记录格式类型（取第一个片段的格式）
			if messageFormat == "text" && chunk.Type != "" {
				messageFormat = chunk.Type
			}

			// 累积完整内容
			fullContent.WriteString(chunk.Content)

			// 通过 WebSocket 推送片段
			sendBotMessageChunk(client, botMsgID, conversationID, botApp.BotUserID, chunk)
			return nil
		},
		OnComplete: func() error {
			// 发送完成标记
			sendBotMessageDone(client, botMsgID, conversationID, botApp.BotUserID, fullContent.String(), messageFormat)
			return nil
		},
		OnError: func(err error) error {
			sendBotMessageError(client, botMsgID, err.Error())
			return nil
		},
	}

	// 6. 发送开始标记
	sendBotMessageStart(client, botMsgID, conversationID, botApp.BotUserID)

	// 7. 执行流式调用
	err = handler.ExecuteStreamCall(ctx)
	if err != nil {
		fmt.Printf("Bot 流式调用失败: %v\n", err)
		sendBotMessageError(client, botMsgID, err.Error())
		return
	}
}

// sendBotMessageStart 发送 Bot 消息开始标记
func sendBotMessageStart(client *ws_conn.Client, msgID, convID, botID string) {
	content := type_struct.WsContent{
		Timestamp: getCurrentTimestamp(),
		MessageID: msgID,
		Data: type_struct.WsData{
			Type:           "bot_streaming_start",
			ConversationID: convID,
			Body:           json.RawMessage(fmt.Sprintf(`{"sender_id":"%s"}`, botID)),
		},
	}

	client.SafeSend(wsCommandConst.CHAT_MESSAGE, content)
}

// sendBotMessageChunk 发送 Bot 消息片段（支持多格式）
func sendBotMessageChunk(client *ws_conn.Client, msgID, convID, botID string, chunk *open_models.BotChunk) {
	bodyMap := map[string]interface{}{
		"sender_id": botID,
		"format":    chunk.Type,
		"content":   chunk.Content,
	}

	if chunk.Metadata != nil && len(chunk.Metadata) > 0 {
		bodyMap["metadata"] = chunk.Metadata
	}

	bodyJSON, _ := json.Marshal(bodyMap)

	content := type_struct.WsContent{
		Timestamp: getCurrentTimestamp(),
		MessageID: msgID,
		Data: type_struct.WsData{
			Type:           "bot_streaming_chunk",
			ConversationID: convID,
			Body:           bodyJSON,
		},
	}

	client.SafeSend(wsCommandConst.CHAT_MESSAGE, content)
}

// sendBotMessageDone 发送 Bot 消息完成标记
func sendBotMessageDone(client *ws_conn.Client, msgID, convID, botID, fullContent string, format string) {
	bodyMap := map[string]interface{}{
		"sender_id": botID,
		"format":    format,
		"content":   fullContent,
	}
	bodyJSON, _ := json.Marshal(bodyMap)

	content := type_struct.WsContent{
		Timestamp: getCurrentTimestamp(),
		MessageID: msgID,
		Data: type_struct.WsData{
			Type:           "bot_streaming_done",
			ConversationID: convID,
			Body:           bodyJSON,
		},
	}

	client.SafeSend(wsCommandConst.CHAT_MESSAGE, content)
}

// sendBotMessageError 发送 Bot 消息错误
func sendBotMessageError(client *ws_conn.Client, msgID, errMsg string) {
	bodyMap := map[string]interface{}{
		"error": errMsg,
	}
	bodyJSON, _ := json.Marshal(bodyMap)

	content := type_struct.WsContent{
		Timestamp: getCurrentTimestamp(),
		MessageID: msgID,
		Data: type_struct.WsData{
			Type: "bot_streaming_error",
			Body: bodyJSON,
		},
	}

	client.SafeSend(wsCommandConst.CHAT_MESSAGE, content)
}

// getCurrentTimestamp 获取当前时间戳（毫秒）
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// generateMessageID 生成消息 ID
func generateMessageID() string {
	return uuid.New().String()
}
