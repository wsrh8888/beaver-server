package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"beaver/app/open/open_models"
	"beaver/common/const/webhookconst"
	"beaver/core/coregorm"

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

// SendCallback 发送 Webhook 回调
func SendCallback(eventType string, payload map[string]interface{}) {
	// 1. 查询订阅该事件的配置
	configs := getWebhookConfigs(eventType)
	if len(configs) == 0 {
		return
	}

	// 2. 构建回调数据
	callbackData := CallbackPayload{
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now().UnixMilli(),
		Payload:   payload,
	}

	// 3. 异步发送回调 (通过 RocketMQ)
	for _, config := range configs {
		go sendCallbackAsync(config, callbackData)
	}
}

// sendCallbackAsync 异步发送回调
func sendCallbackAsync(config open_models.OpenWebhookConfig, data CallbackPayload) {
	// 序列化数据
	body, err := json.Marshal(data)
	if err != nil {
		logx.Errorf("Webhook 序列化失败: %v", err)
		return
	}

	// 生成签名
	signature := generateSignature(body, config.Secret)

	// 发送 HTTP POST
	success := false
	for retry := 0; retry < config.RetryCount; retry++ {
		if retry > 0 {
			// 指数退避
			time.Sleep(time.Duration(retry*retry) * time.Second)
		}

		if sendHTTPRequest(config.TargetURL, body, signature, config.Timeout) {
			success = true
			logx.Infof("Webhook 回调成功: EventID=%s, URL=%s", data.EventID, config.TargetURL)
			break
		}
		logx.Warnf("Webhook 回调失败 (重试 %d/%d): EventID=%s", retry+1, config.RetryCount, data.EventID)
	}

	// 记录日志
	saveWebhookLog(config, data, success)
}

// sendHTTPRequest 发送 HTTP 请求
func sendHTTPRequest(url string, body []byte, signature string, timeout int) bool {
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
	req.Header.Set("X-Webhook-Event-ID", "") // 从 data 获取

	resp, err := client.Do(req)
	if err != nil {
		logx.Errorf("HTTP 请求失败: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// generateSignature 生成签名
func generateSignature(body []byte, secret string) string {
	if secret == "" {
		return ""
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

// getWebhookConfigs 查询 Webhook 配置
func getWebhookConfigs(eventType string) []open_models.OpenWebhookConfig {
	db := coregorm.GetDB()
	var configs []open_models.OpenWebhookConfig
	db.Where("event_type = ? AND status = ?", eventType, 1).Find(&configs)
	return configs
}

// saveWebhookLog 保存 Webhook 日志
func saveWebhookLog(config open_models.OpenWebhookConfig, data CallbackPayload, success bool) {
	db := coregorm.GetDB()

	status := 0
	if success {
		status = 1
	}

	log := open_models.OpenWebhookLog{
		ConfigID:  fmt.Sprintf("%d", config.ID),
		AppID:     config.AppID,
		EventType: data.EventType,
		Status:    status,
	}

	// 简化处理，实际应该从 MQ 或上下文获取完整信息
	db.Create(&log)
}

// TriggerBotCallback 触发 Bot 消息回调
func TriggerBotCallback(botAppID string, sender string, content string, conversationType int, groupID string, msgID string) {
	payload := map[string]interface{}{
		"bot_app_id":        botAppID,
		"sender":            sender,
		"content":           content,
		"conversation_type": conversationType,
		"group_id":          groupID,
		"msg_id":            msgID,
		"timestamp":         time.Now().UnixMilli(),
	}

	SendCallback(webhookconst.EventBotMessageReceive, payload)
}
