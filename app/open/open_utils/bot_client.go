package open_utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// BotPushRequest Bot 主动推送消息请求（大厂标准）
type BotPushRequest struct {
	AppID          string                 `json:"app_id"`                   // 应用 ID
	ConversationID string                 `json:"conversation_id"`          // 会话 ID（私聊或群聊）
	MsgType        string                 `json:"msg_type"`                 // 消息类型：text/markdown/richtext
	Content        string                 `json:"content"`                  // 消息内容
	Metadata       map[string]interface{} `json:"metadata,omitempty"`       // 元数据
	IdempotentKey  string                 `json:"idempotent_key,omitempty"` // 幂等键（防止重复）
}

// BotPushResponse Bot 推送响应
type BotPushResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
}

// BotClient Bot SDK 客户端（提供给第三方开发者）
type BotClient struct {
	BaseURL    string
	AppID      string
	AppSecret  string
	HTTPClient *http.Client
}

// NewBotClient 创建 Bot 客户端
func NewBotClient(baseURL, appID, appSecret string) *BotClient {
	return &BotClient{
		BaseURL:   baseURL,
		AppID:     appID,
		AppSecret: appSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendMessage Bot 主动发送消息（对标飞书/钉钉 Bot API）
func (c *BotClient) SendMessage(ctx context.Context, req *BotPushRequest) (*BotPushResponse, error) {
	req.AppID = c.AppID

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST",
		c.BaseURL+"/api/open/v1/bot/message/send", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AppSecret))

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result BotPushResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if !result.Success {
		return &result, fmt.Errorf("Bot 推送失败: %s", result.Error)
	}

	return &result, nil
}

// SendTextMessage 发送文本消息
func (c *BotClient) SendTextMessage(ctx context.Context, conversationID, content string) (*BotPushResponse, error) {
	return c.SendMessage(ctx, &BotPushRequest{
		ConversationID: conversationID,
		MsgType:        "text",
		Content:        content,
	})
}

// SendMarkdownMessage 发送 Markdown 消息
func (c *BotClient) SendMarkdownMessage(ctx context.Context, conversationID, content string) (*BotPushResponse, error) {
	return c.SendMessage(ctx, &BotPushRequest{
		ConversationID: conversationID,
		MsgType:        "markdown",
		Content:        content,
	})
}

// SendStreamMessage 发送流式消息（SSE）
func (c *BotClient) SendStreamMessage(ctx context.Context, conversationID string, chunks <-chan string) error {
	// 创建流式请求
	req, err := http.NewRequestWithContext(ctx, "POST",
		c.BaseURL+"/api/open/v1/bot/message/stream", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/event-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AppSecret))

	// TODO: 实现 SSE 客户端发送
	// 这里简化处理，实际应该用专门的 SSE 客户端库

	return nil
}
