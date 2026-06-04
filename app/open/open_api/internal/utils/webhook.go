package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"beaver/app/open/open_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type urlVerificationPayload struct {
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Token     string `json:"token,omitempty"`
}

type urlVerificationResponse struct {
	Challenge string `json:"challenge"`
}

type platformEventPayload struct {
	Type      string                 `json:"type"`
	EventID   string                 `json:"event_id"`
	EventType string                 `json:"event_type"`
	Timestamp int64                  `json:"timestamp"`
	Token     string                 `json:"token,omitempty"`
	Event     map[string]interface{} `json:"event"`
}

func VerifyWebhookURL(targetURL, secret string, timeoutSec int) error {
	if targetURL == "" {
		return errors.New("targetUrl 不能为空")
	}
	if timeoutSec <= 0 {
		timeoutSec = 5
	}

	challengeBytes := make([]byte, 16)
	if _, err := rand.Read(challengeBytes); err != nil {
		return err
	}
	challenge := hex.EncodeToString(challengeBytes)

	payload := urlVerificationPayload{Type: "url_verification", Challenge: challenge}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		req.Header.Set("X-Webhook-Signature", signBody(body, secret))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Challenge 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Challenge 验证失败: HTTP %d", resp.StatusCode)
	}

	var result urlVerificationResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return errors.New("Challenge 响应不是合法 JSON")
	}
	if result.Challenge != challenge {
		return errors.New("Challenge 响应不匹配")
	}
	return nil
}

func PushPlatformEvent(sub open_models.OpenAppEventSubscription, eventType string, event map[string]interface{}) error {
	if sub.Status != 1 || sub.VerifyStatus != 1 {
		return nil
	}

	payload := platformEventPayload{
		Type:      "event_callback",
		EventID:   uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now().Unix(),
		Event:     event,
	}
	if sub.Secret != "" {
		payload.Token = signEventToken(sub.Secret, payload.EventID, payload.Timestamp)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	timeout := sub.Timeout
	if timeout <= 0 {
		timeout = 5
	}
	retryCount := sub.RetryCount
	if retryCount <= 0 {
		retryCount = 3
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	for i := 0; i < retryCount; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i*i) * time.Second)
		}
		req, err := http.NewRequest(http.MethodPost, sub.CallbackURL, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		if sub.Secret != "" {
			req.Header.Set("X-Webhook-Signature", signBody(body, sub.Secret))
			req.Header.Set("X-Webhook-Event-ID", payload.EventID)
			req.Header.Set("X-Webhook-Timestamp", strconv.FormatInt(payload.Timestamp, 10))
		}

		resp, err := client.Do(req)
		if err != nil {
			logx.Errorf("[webhook] push failed: app=%s event=%s err=%v", sub.AppID, eventType, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		logx.Errorf("[webhook] push bad status: app=%s event=%s code=%d", sub.AppID, eventType, resp.StatusCode)
	}
	return fmt.Errorf("webhook push failed after retries")
}

func signBody(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func signEventToken(secret, eventID string, timestamp int64) string {
	raw := fmt.Sprintf("%s:%d", eventID, timestamp)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}
