package corewebhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

// CallbackPayload Webhook 回调载荷
type CallbackPayload struct {
	EventID   string                 `json:"eventId"`   // 事件唯一ID
	EventType string                 `json:"eventType"` // 事件类型
	Timestamp int64                  `json:"timestamp"` // 时间戳
	Payload   map[string]interface{} `json:"payload"`   // 业务数据
}

// Config Webhook 配置
type Config struct {
	Timeout    int
	RetryCount int
}

// LogWriter 跨域写 Webhook 日志（由 open_rpc 等实现）
type LogWriter interface {
	SaveWebhookLog(ctx context.Context, configID, appID, eventType string, success bool)
}

// WebhookSender Webhook 发送器
type WebhookSender struct {
	logWriter LogWriter
	config    Config
}

// NewWebhookSender 创建发送器，logWriter 可为 nil（不记日志）
func NewWebhookSender(conf Config, logWriter LogWriter) *WebhookSender {
	if conf.Timeout == 0 {
		conf.Timeout = 10
	}
	return &WebhookSender{
		logWriter: logWriter,
		config:    conf,
	}
}

// Send 发送 Webhook 回调
func (s *WebhookSender) Send(eventType string, payload map[string]interface{}, configs []WebhookTargetConfig) {
	if len(configs) == 0 {
		return
	}

	// 构建基础数据
	data := CallbackPayload{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}

	// 异步发送 (大厂建议：这里以后可以改造成发送给 RocketMQ)
	for _, config := range configs {
		go s.sendAsync(config, data)
	}
}

// WebhookTargetConfig 目标配置（从业务层传入，解耦核心库与特定的 OpenWebhookConfig 模型）
type WebhookTargetConfig struct {
	ID         uint
	AppID      string
	TargetURL  string
	Secret     string
	RetryCount int
	Timeout    int
}

func (s *WebhookSender) sendAsync(config WebhookTargetConfig, data CallbackPayload) {
	body, err := json.Marshal(data)
	if err != nil {
		logx.Errorf("Webhook 序列化失败: %v", err)
		return
	}

	signature := s.generateSignature(body, config.Secret)

	retryCount := config.RetryCount
	if retryCount == 0 {
		retryCount = s.config.RetryCount
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = s.config.Timeout
	}

	success := false
	for retry := 0; retry < retryCount; retry++ {
		if retry > 0 {
			// 指数退避
			time.Sleep(time.Duration(retry*retry) * time.Second)
		}

		if s.sendHTTPRequest(config.TargetURL, body, signature, data.EventID, data.Timestamp, timeout) {
			success = true
			logx.Infof("Webhook 回调成功: EventID=%s, URL=%s", data.EventID, config.TargetURL)
			break
		}
		logx.Errorf("Webhook 回调失败 (重试 %d/%d): EventID=%s", retry+1, retryCount, data.EventID)
	}

	s.saveLog(config, data, success)
}

func (s *WebhookSender) sendHTTPRequest(url string, body []byte, signature, eventID string, timestamp int64, timeout int) bool {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		logx.Errorf("创建 HTTP 请求失败: %v", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Event-ID", eventID)
	req.Header.Set("X-Webhook-Timestamp", strconv.FormatInt(timestamp, 10))

	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("HTTP 请求失败: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func (s *WebhookSender) generateSignature(body []byte, secret string) string {
	if secret == "" {
		return ""
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *WebhookSender) saveLog(config WebhookTargetConfig, data CallbackPayload, success bool) {
	if s.logWriter == nil {
		return
	}
	s.logWriter.SaveWebhookLog(context.Background(), fmt.Sprintf("%d", config.ID), config.AppID, data.EventType, success)
}
