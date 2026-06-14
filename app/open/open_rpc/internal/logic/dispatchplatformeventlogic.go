package logic

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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"

	"gorm.io/gorm"
)

type DispatchPlatformEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchPlatformEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchPlatformEventLogic {
	return &DispatchPlatformEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DispatchPlatformEventLogic) DispatchPlatformEvent(in *open_rpc.DispatchPlatformEventReq) (*open_rpc.DispatchPlatformEventRes, error) {
	if in.AppId == "" || in.EventType == "" || in.EventJson == "" {
		return &open_rpc.DispatchPlatformEventRes{Dispatched: false}, nil
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", in.AppId).First(&app).Error; err != nil {
		return &open_rpc.DispatchPlatformEventRes{Dispatched: false}, nil
	}
	if app.EnableWebhook != 1 || app.EnableRobot != 1 {
		return &open_rpc.DispatchPlatformEventRes{Dispatched: false}, nil
	}

	var event map[string]interface{}
	if err := json.Unmarshal([]byte(in.EventJson), &event); err != nil {
		l.Errorf("DispatchPlatformEvent 事件 JSON 无效: %v", err)
		return &open_rpc.DispatchPlatformEventRes{Dispatched: false}, nil
	}

	subs, err := listActiveSubscriptions(l.svcCtx.DB, in.AppId, in.EventType)
	if err != nil || len(subs) == 0 {
		return &open_rpc.DispatchPlatformEventRes{Dispatched: false}, nil
	}

	dispatched := false
	for _, sub := range subs {
		eventID, httpStatus, latencyMs, retries, pushErr := pushPlatformEvent(sub, in.EventType, event)
		status := 1
		errMsg := ""
		if pushErr != nil {
			status = 0
			errMsg = pushErr.Error()
		} else {
			dispatched = true
		}
		_ = l.svcCtx.DB.Create(&open_models.OpenWebhookLog{
			SubscriptionID: sub.ID,
			AppID:          in.AppId,
			EventID:        eventID,
			EventType:      in.EventType,
			TargetURL:      sub.CallbackURL,
			HTTPStatus:     httpStatus,
			LatencyMs:      latencyMs,
			RetryCount:     retries,
			Status:         status,
			ErrorMessage:   errMsg,
		}).Error
	}

	return &open_rpc.DispatchPlatformEventRes{Dispatched: dispatched}, nil
}

func listActiveSubscriptions(db *gorm.DB, appID, eventType string) ([]open_models.OpenAppEventSubscription, error) {
	var subs []open_models.OpenAppEventSubscription
	err := db.Where("app_id = ? AND event_type = ? AND status = 1 AND verify_status = 1", appID, eventType).
		Find(&subs).Error
	return subs, err
}

func pushPlatformEvent(sub open_models.OpenAppEventSubscription, eventType string, event map[string]interface{}) (eventID string, httpStatus int, latencyMs int64, retryCount int, err error) {
	if sub.Status != 1 || sub.VerifyStatus != 1 {
		return "", 0, 0, 0, nil
	}

	eventID = uuid.New().String()
	timestamp := time.Now().Unix()
	payload := map[string]interface{}{
		"type":       "event_callback",
		"event_id":   eventID,
		"event_type": eventType,
		"timestamp":  timestamp,
		"event":      event,
	}
	if sub.Secret != "" {
		payload["token"] = signEventToken(sub.Secret, eventID, timestamp)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return eventID, 0, 0, 0, err
	}

	timeout := sub.Timeout
	if timeout <= 0 {
		timeout = 5
	}
	maxRetry := sub.RetryCount
	if maxRetry <= 0 {
		maxRetry = 3
	}

	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	var lastStatus int
	for i := 0; i < maxRetry; i++ {
		retryCount = i + 1
		if i > 0 {
			time.Sleep(time.Duration(i*i) * time.Second)
		}
		start := time.Now()
		req, err := http.NewRequest(http.MethodPost, sub.CallbackURL, bytes.NewReader(body))
		if err != nil {
			return eventID, 0, 0, retryCount, err
		}
		req.Header.Set("Content-Type", "application/json")
		if sub.Secret != "" {
			req.Header.Set("X-Webhook-Signature", signBody(body, sub.Secret))
			req.Header.Set("X-Webhook-Event-ID", eventID)
			req.Header.Set("X-Webhook-Timestamp", strconv.FormatInt(timestamp, 10))
		}

		resp, err := client.Do(req)
		latencyMs = time.Since(start).Milliseconds()
		if err != nil {
			logx.Errorf("[open_rpc] webhook push failed: app=%s event=%s err=%v", sub.AppID, eventType, err)
			continue
		}
		lastStatus = resp.StatusCode
		resp.Body.Close()
		httpStatus = lastStatus
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return eventID, httpStatus, latencyMs, retryCount, nil
		}
		logx.Errorf("[open_rpc] webhook push bad status: app=%s event=%s code=%d", sub.AppID, eventType, resp.StatusCode)
	}
	return eventID, httpStatus, latencyMs, retryCount, fmt.Errorf("webhook push failed after retries")
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
