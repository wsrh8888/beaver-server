package event

import (
	"context"
	"errors"
	"strconv"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEventSubscriptionLogic {
	return &UpdateEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEventSubscriptionLogic) UpdateEventSubscription(req *types.UpdateEventSubscriptionReq) (resp *types.UpdateEventSubscriptionRes, err error) {
	subID, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil || subID == 0 {
		return nil, errors.New("订阅 ID 无效")
	}

	var sub open_models.OpenAppEventSubscription
	if err := l.svcCtx.DB.Where("id = ?", subID).First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("订阅不存在")
		}
		return nil, errors.New("查询订阅失败")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", sub.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	updates := map[string]interface{}{}
	if req.TargetURL != "" {
		updates["callback_url"] = req.TargetURL
	}
	if req.Secret != "" {
		updates["secret"] = req.Secret
	}
	if req.Status != nil {
		if *req.Status != 0 && *req.Status != 1 {
			return nil, errors.New("status 只能为 0 或 1")
		}
		updates["status"] = *req.Status
	}
	if req.RetryCount != nil {
		updates["retry_count"] = *req.RetryCount
	}
	if req.Timeout != nil {
		updates["timeout"] = *req.Timeout
	}

	if len(updates) == 0 {
		return &types.UpdateEventSubscriptionRes{}, nil
	}

	targetURL := sub.CallbackURL
	if v, ok := updates["callback_url"].(string); ok {
		targetURL = v
	}
	secret := sub.Secret
	if v, ok := updates["secret"].(string); ok {
		secret = v
	}
	timeout := sub.Timeout
	if v, ok := updates["timeout"].(int); ok {
		timeout = v
	}

	if req.TargetURL != "" {
		verifyErr := verifyWebhookChallenge(targetURL, secret, timeout)
		now := time.Now()
		if verifyErr != nil {
			updates["verify_status"] = 2
			updates["last_error"] = verifyErr.Error()
		} else {
			updates["verify_status"] = 1
			updates["last_verified_at"] = &now
			updates["last_error"] = ""
		}
	}

	if err := l.svcCtx.DB.Model(&sub).Updates(updates).Error; err != nil {
		return nil, errors.New("更新订阅失败")
	}

	return &types.UpdateEventSubscriptionRes{}, nil
}
