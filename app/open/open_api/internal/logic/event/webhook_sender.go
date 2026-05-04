package event

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"beaver-server/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// WebhookEvent 事件数据结构
type WebhookEvent struct {
	EventID   string      `json:"event_id"`
	EventType string      `json:"event_type"`
	Timestamp int64       `json:"timestamp"`
	AppID     string      `json:"app_id"`
	Data      interface{} `json:"data"`
}

// GenerateSignature 生成 HMAC-SHA256 签名
func GenerateSignature(payload string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// SendWebhook 发送 Webhook 事件(带重试机制)
func SendWebhook(db *gorm.DB, config models.OpenWebhookConfig, eventData interface{}) error {
	// 1. 构建事件数据
	event := WebhookEvent{
		EventID:   fmt.Sprintf("evt_%d", time.Now().UnixNano()),
		EventType: config.EventType,
		Timestamp: time.Now().Unix(),
		AppID:     config.AppID,
		Data:      eventData,
	}

	// 2. 序列化 payload
	payloadBytes, err := json.Marshal(event)
	if err != nil {
		logx.Errorf("序列化 Webhook 事件失败: %v", err)
		return err
	}
	payloadStr := string(payloadBytes)

	// 3. 生成签名
	signature := GenerateSignature(payloadStr, config.Secret)

	// 4. 发送请求(带重试)
	maxRetries := config.RetryCount
	timeout := time.Duration(config.Timeout) * time.Second

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// 指数退避: 1s, 2s, 4s, 8s...
			waitTime := time.Duration(1<<uint(attempt-1)) * time.Second
			logx.Infof("Webhook 重试 %d/%d, 等待 %v", attempt, maxRetries, waitTime)
			time.Sleep(waitTime)
		}

		err = sendWithRetry(config.TargetURL, payloadStr, signature, timeout)
		if err == nil {
			// 成功,记录日志
			logWebhook(db, config, payloadStr, 200, attempt, 1)
			logx.Infof("Webhook 推送成功: event_type=%s, url=%s", config.EventType, config.TargetURL)
			return nil
		}

		lastErr = err
		logx.Errorf("Webhook 推送失败 (尝试 %d/%d): %v", attempt+1, maxRetries+1, err)
	}

	// 全部失败,记录日志
	logWebhook(db, config, payloadStr, 0, maxRetries, 2)
	logx.Errorf("Webhook 推送最终失败: event_type=%s, url=%s", config.EventType, config.TargetURL)
	return lastErr
}

// sendWithRetry 单次 HTTP 请求
func sendWithRetry(url string, payload string, signature string, timeout time.Duration) error {
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Beaver-Signature", signature)
	req.Header.Set("X-Beaver-Event-Type", req.Header.Get("X-Beaver-Event-Type"))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应体(用于调试)
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// logWebhook 记录 Webhook 日志
func logWebhook(db *gorm.DB, config models.OpenWebhookConfig, payload string, responseCode int, retryCount int, status int) {
	log := models.OpenWebhookLog{
		AppID:        config.AppID,
		EventType:    config.EventType,
		Payload:      payload,
		ResponseCode: responseCode,
		RetryCount:   retryCount,
		Status:       status,
	}
	db.Create(&log)
}

// TriggerMessageEvent 触发消息事件
// 注意:此函数是异步的,不会阻塞主流程,失败也不会影响业务
func TriggerMessageEvent(db *gorm.DB, appID string, messageData interface{}) {
	// 查询该应用的消息事件配置
	var config models.OpenWebhookConfig
	err := db.Where("app_id = ? AND event_type = ? AND status = ?", appID, "message.receive", 1).First(&config).Error
	if err != nil {
		// 未配置或已禁用,静默返回(正常情况)
		return
	}

	// 异步推送,不阻塞主流程
	go func() {
		err := SendWebhook(db, config, messageData)
		if err != nil {
			// 推送失败,记录详细日志供排查
			logx.Errorf("[Webhook] 消息事件推送失败: app_id=%s, error=%v", appID, err)
			// TODO: 可以加入重试队列或发送告警
		}
	}()
}

// TriggerGroupMemberEvent 触发群成员变更事件
// 注意:此函数是异步的,不会阻塞主流程,失败也不会影响业务
func TriggerGroupMemberEvent(db *gorm.DB, appID string, groupData interface{}) {
	var config models.OpenWebhookConfig
	err := db.Where("app_id = ? AND event_type = ? AND status = ?", appID, "group.member.change", 1).First(&config).Error
	if err != nil {
		// 未配置或已禁用,静默返回(正常情况)
		return
	}

	// 异步推送,不阻塞主流程
	go func() {
		err := SendWebhook(db, config, groupData)
		if err != nil {
			// 推送失败,记录详细日志供排查
			logx.Errorf("[Webhook] 群成员事件推送失败: app_id=%s, error=%v", appID, err)
			// TODO: 可以加入重试队列或发送告警
		}
	}()
}
