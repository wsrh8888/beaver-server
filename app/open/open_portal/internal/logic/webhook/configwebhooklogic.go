package webhook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	models "beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 配置 Webhook
func NewConfigWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigWebhookLogic {
	return &ConfigWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigWebhookLogic) ConfigWebhook(req *types.ConfigWebhookReq) (resp *types.ConfigWebhookRes, err error) {
	// 1. 验证应用是否存在
	var app models.OpenApp
	err = l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error
	if err != nil {
		return nil, errors.New("应用不存在")
	}

	// 2. 生成签名密钥(如果未提供)
	secret := req.Secret
	if secret == "" {
		secretBytes := make([]byte, 32)
		rand.Read(secretBytes)
		secret = hex.EncodeToString(secretBytes)
	}

	// 3. 设置默认值
	retryCount := req.RetryCount
	if retryCount == 0 {
		retryCount = 3
	}
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 5
	}

	// 4. 检查是否已存在相同事件类型的配置
	var existingConfig models.OpenWebhookConfig
	err = l.svcCtx.DB.Where("app_id = ? AND event_type = ?", req.AppID, req.EventType).First(&existingConfig).Error
	if err == nil {
		// 更新现有配置
		existingConfig.TargetURL = req.TargetURL
		existingConfig.Secret = secret
		if req.Events != "" {
			existingConfig.Events = req.Events
		}
		existingConfig.RetryCount = retryCount
		existingConfig.Timeout = timeout
		existingConfig.Status = 1
		l.svcCtx.DB.Save(&existingConfig)

		logx.Infof("更新 Webhook 配置: app_id=%s, event_type=%s", req.AppID, req.EventType)

		return &types.ConfigWebhookRes{
			ConfigID: fmt.Sprintf("%d", existingConfig.ID),
		}, nil
	}

	// 5. 创建新配置
	newConfig := models.OpenWebhookConfig{
		AppID:      req.AppID,
		EventType:  req.EventType,
		TargetURL:  req.TargetURL,
		Secret:     secret,
		Events:     req.Events,
		Status:     1,
		RetryCount: retryCount,
		Timeout:    timeout,
	}

	err = l.svcCtx.DB.Create(&newConfig).Error
	if err != nil {
		return nil, errors.New("创建 Webhook 配置失败")
	}

	logx.Infof("创建 Webhook 配置: app_id=%s, event_type=%s", req.AppID, req.EventType)

	return &types.ConfigWebhookRes{
		ConfigID: fmt.Sprintf("%d", newConfig.ID),
	}, nil
}
