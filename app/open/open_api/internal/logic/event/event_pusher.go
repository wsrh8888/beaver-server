package event

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"beaver-server/app/open/open_models"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// EventPayload 事件数据结构
type EventPayload struct {
	EventID   string      `json:"eventId"`
	EventType string      `json:"eventType"`
	Timestamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// PushEvent 推送事件到开发者服务器
func PushEvent(db *gorm.DB, appID, eventType string, data interface{}) error {
	// 1. 查询该应用的事件订阅配置
	var subscription open_models.OpenEventSubscription
	if err := db.Where("app_id = ? AND event_type = ? AND status = ?", appID, eventType, 1).First(&subscription).Error; err != nil {
		// 没有订阅，直接返回
		return nil
	}

	// 2. 构建事件数据
	eventID := fmt.Sprintf("evt_%d", time.Now().UnixNano())
	payload := EventPayload{
		EventID:   eventID,
		EventType: eventType,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logx.Errorf("序列化事件数据失败: %v", err)
		return err
	}

	// 3. 创建事件日志
	eventLog := open_models.OpenEventLog{
		EventID:    eventID,
		AppID:      appID,
		EventType:  eventType,
		Payload:    string(payloadJSON),
		TargetURL:  subscription.TargetURL,
		RetryCount: 0,
		Status:     0, // 待推送
		CreatedAt:  time.Now().Unix(),
	}

	if err := db.Create(&eventLog).Error; err != nil {
		logx.Errorf("创建事件日志失败: %v", err)
		return err
	}

	// 4. 发送 HTTP 请求
	startTime := time.Now()
	
	req, err := http.NewRequest("POST", subscription.TargetURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		logx.Errorf("创建 HTTP 请求失败: %v", err)
		updateEventLog(db, eventLog.ID, 0, 0, 2, err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	
	// 添加签名头（如果配置了 secret）
	if subscription.Secret != "" {
		signature := generateSignature(string(payloadJSON), subscription.Secret)
		req.Header.Set("X-Beaver-Signature", signature)
		req.Header.Set("X-Beaver-Timestamp", fmt.Sprintf("%d", payload.Timestamp))
	}

	client := &http.Client{
		Timeout: time.Duration(subscription.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	responseTime := int(time.Since(startTime).Milliseconds())

	if err != nil {
		logx.Errorf("推送事件失败: %v", err)
		updateEventLog(db, eventLog.ID, 0, responseTime, 2, err.Error())
		return err
	}
	defer resp.Body.Close()

	// 5. 更新事件日志
	status := 1 // 成功
	errorMsg := ""
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		status = 2 // 失败
		errorMsg = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	updateEventLog(db, eventLog.ID, resp.StatusCode, responseTime, status, errorMsg)

	if status == 2 {
		return fmt.Errorf("推送事件失败: %s", errorMsg)
	}

	return nil
}

// generateSignature 生成 HMAC-SHA256 签名
func generateSignature(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// updateEventLog 更新事件日志
func updateEventLog(db *gorm.DB, logID uint64, responseCode, responseTime, status int, errorMsg string) {
	updates := map[string]interface{}{
		"response_code": responseCode,
		"response_time": responseTime,
		"status":        status,
		"error_msg":     errorMsg,
	}

	if err := db.Model(&open_models.OpenEventLog{}).Where("id = ?", logID).Updates(updates).Error; err != nil {
		logx.Errorf("更新事件日志失败: %v", err)
	}
}
