package event

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/openevent"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CreateEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEventSubscriptionLogic {
	return &CreateEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEventSubscriptionLogic) CreateEventSubscription(req *types.CreateEventSubscriptionReq) (resp *types.CreateEventSubscriptionRes, err error) {
	if req.AppID == "" || req.EventType == "" || req.TargetURL == "" {
		return nil, errors.New("appId、eventType 和 targetUrl 不能为空")
	}
	if err := openevent.ValidateRobotEventType(req.EventType); err != nil {
		return nil, err
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}
	if app.EnableWebhook != 1 {
		return nil, errors.New("应用未启用 Webhook 能力")
	}

	retryCount := req.RetryCount
	if retryCount <= 0 {
		retryCount = 3
	}
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 5
	}

	var existing open_models.OpenAppEventSubscription
	err = l.svcCtx.DB.Where("app_id = ? AND event_type = ?", req.AppID, req.EventType).First(&existing).Error
	if err == nil {
		return nil, errors.New("该事件类型已存在订阅，请使用更新接口")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("查询订阅失败")
	}

	sub := open_models.OpenAppEventSubscription{
		AppID:       req.AppID,
		EventType:   req.EventType,
		CallbackURL: req.TargetURL,
		Secret:      req.Secret,
		Status:      1,
		RetryCount:  retryCount,
		Timeout:     timeout,
	}

	verifyErr := verifyWebhookChallenge(req.TargetURL, req.Secret, timeout)
	now := time.Now()
	if verifyErr != nil {
		sub.VerifyStatus = 2
		sub.LastError = verifyErr.Error()
	} else {
		sub.VerifyStatus = 1
		sub.LastVerifiedAt = &now
	}

	if err := l.svcCtx.DB.Create(&sub).Error; err != nil {
		return nil, errors.New("创建事件订阅失败")
	}

	status := sub.Status
	if sub.VerifyStatus != 1 {
		status = 0
	}

	lastVerifiedAt := int64(0)
	if sub.LastVerifiedAt != nil {
		lastVerifiedAt = sub.LastVerifiedAt.Unix()
	}
	return &types.CreateEventSubscriptionRes{
		Subscription: types.CreateEventSubscriptionResSubscription{
			ID:             fmt.Sprintf("%d", sub.ID),
			AppID:          sub.AppID,
			EventType:      sub.EventType,
			TargetURL:      sub.CallbackURL,
			Secret:         sub.Secret,
			Status:         status,
			VerifyStatus:   sub.VerifyStatus,
			LastError:      sub.LastError,
			LastVerifiedAt: lastVerifiedAt,
			RetryCount:     sub.RetryCount,
			Timeout:        sub.Timeout,
			CreatedAt:      sub.CreatedAt.Unix(),
			UpdatedAt:      sub.UpdatedAt.Unix(),
		},
	}, nil
}

func verifyWebhookChallenge(targetURL, secret string, timeoutSec int) error {
	challengeBytes := make([]byte, 16)
	if _, err := rand.Read(challengeBytes); err != nil {
		return err
	}
	challenge := hex.EncodeToString(challengeBytes)

	body, err := json.Marshal(map[string]string{
		"type":      "url_verification",
		"challenge": challenge,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if secret != "" {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		req.Header.Set("X-Webhook-Signature", hex.EncodeToString(mac.Sum(nil)))
	}

	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Challenge 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Challenge 验证失败: HTTP %d", resp.StatusCode)
	}

	var result struct {
		Challenge string `json:"challenge"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return errors.New("Challenge 响应不是合法 JSON")
	}
	if result.Challenge != challenge {
		return errors.New("Challenge 响应不匹配")
	}
	return nil
}
