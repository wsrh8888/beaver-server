package event

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEventSubscriptionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建事件订阅
func NewCreateEventSubscriptionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEventSubscriptionLogic {
	return &CreateEventSubscriptionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateEventSubscriptionLogic) CreateEventSubscription(req *types.CreateEventSubscriptionReq) (resp *types.CreateEventSubscriptionRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// 创建事件订阅记录
	subscription := models.OpenAppEventSubscription{
		AppID:       req.AppID,
		EventType:   req.EventType,
		CallbackURL: req.TargetURL,
		Secret:      req.Secret,
		Status:      1,
		RetryCount:  req.RetryCount,
		Timeout:     req.Timeout,
	}

	if req.RetryCount == 0 {
		subscription.RetryCount = 3
	}
	if req.Timeout == 0 {
		subscription.Timeout = 5
	}

	if err := l.svcCtx.DB.Create(&subscription).Error; err != nil {
		return nil, errors.New("创建事件订阅失败")
	}

	// 2. 返回结果
	return &types.CreateEventSubscriptionRes{
		Subscription: types.EventSubscriptionInfo{
			ID:         fmt.Sprintf("%d", subscription.ID),
			AppID:      subscription.AppID,
			EventType:  subscription.EventType,
			TargetURL:  subscription.CallbackURL,
			Secret:     subscription.Secret,
			Status:     subscription.Status,
			RetryCount: subscription.RetryCount,
			Timeout:    subscription.Timeout,
			CreatedAt:  time.Now().Unix(),
			UpdatedAt:  time.Now().Unix(),
		},
	}, nil
}
