package event

import (
	"context"
	"fmt"
	"time"

	"beaver-server/app/open/open_models"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigEventSubscriptionLogic {
	return &ConfigEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ConfigEventSubscription 配置事件订阅
func (l *ConfigEventSubscriptionLogic) ConfigEventSubscription(req *types.ConfigEventSubscriptionReq) (*types.ConfigEventSubscriptionRes, error) {
	// 1. 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, fmt.Errorf("应用不存在")
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, fmt.Errorf("应用已禁用")
	}

	// 3. 检查是否已存在相同的事件订阅
	var existing open_models.OpenEventSubscription
	err := l.svcCtx.DB.Where("app_id = ? AND event_type = ?", req.AppID, req.EventType).First(&existing).Error
	if err == nil {
		// 更新现有订阅
		existing.TargetURL = req.TargetURL
		existing.Secret = req.Secret
		existing.RetryCount = req.RetryCount
		existing.Timeout = req.Timeout
		existing.Status = 1
		existing.UpdatedAt = time.Now().Unix()
		
		if err := l.svcCtx.DB.Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("更新订阅失败")
		}
		
		return &types.ConfigEventSubscriptionRes{
			SubscriptionID: existing.ID,
		}, nil
	}

	// 4. 创建新订阅
	subscription := open_models.OpenEventSubscription{
		AppID:      req.AppID,
		EventType:  req.EventType,
		TargetURL:  req.TargetURL,
		Secret:     req.Secret,
		RetryCount: req.RetryCount,
		Timeout:    req.Timeout,
		Status:     1,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	if err := l.svcCtx.DB.Create(&subscription).Error; err != nil {
		return nil, fmt.Errorf("创建订阅失败")
	}

	return &types.ConfigEventSubscriptionRes{
		SubscriptionID: subscription.ID,
	}, nil
}
