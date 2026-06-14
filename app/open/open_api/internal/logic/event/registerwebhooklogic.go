package event

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"
	"beaver/app/open/open_models"
	"beaver/app/open/openevent"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterWebhookLogic {
	return &RegisterWebhookLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RegisterWebhookLogic) RegisterWebhook(req *types.RegisterWebhookReq, authorization string) (resp *types.RegisterWebhookRes, err error) {
	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}
	app, err := utils.LoadAppByID(l.svcCtx.DB, token.AppID)
	if err != nil {
		return nil, err
	}
	if err := utils.RequireAppCapability(app, false, true); err != nil {
		return nil, err
	}
	if req.EventType == "" || req.TargetURL == "" {
		return nil, errors.New("eventType 和 targetUrl 不能为空")
	}
	if err := openevent.ValidateRobotEventType(req.EventType); err != nil {
		return nil, err
	}

	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 5
	}
	retryCount := req.RetryCount
	if retryCount <= 0 {
		retryCount = 3
	}

	verifyErr := utils.VerifyWebhookURL(req.TargetURL, req.Secret, timeout)
	now := time.Now()
	sub := open_models.OpenAppEventSubscription{
		AppID:       token.AppID,
		EventType:   req.EventType,
		CallbackURL: req.TargetURL,
		Secret:      req.Secret,
		Status:      1,
		RetryCount:  retryCount,
		Timeout:     timeout,
	}
	if verifyErr != nil {
		sub.VerifyStatus = 2
		sub.LastError = verifyErr.Error()
	} else {
		sub.VerifyStatus = 1
		sub.LastVerifiedAt = &now
	}

	var existing open_models.OpenAppEventSubscription
	err = l.svcCtx.DB.Where("app_id = ? AND event_type = ?", token.AppID, req.EventType).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := l.svcCtx.DB.Create(&sub).Error; err != nil {
			return nil, errors.New("注册 Webhook 失败")
		}
	} else if err != nil {
		return nil, errors.New("查询订阅失败")
	} else {
		sub.ID = existing.ID
		sub.CreatedAt = existing.CreatedAt
		if err := l.svcCtx.DB.Save(&sub).Error; err != nil {
			return nil, errors.New("更新 Webhook 失败")
		}
	}

	status := "active"
	if verifyErr != nil {
		return &types.RegisterWebhookRes{
			SubscriptionID: uint64(sub.ID),
			Status:         "pending",
		}, nil
	}

	return &types.RegisterWebhookRes{
		SubscriptionID: uint64(sub.ID),
		Status:         status,
	}, nil
}
