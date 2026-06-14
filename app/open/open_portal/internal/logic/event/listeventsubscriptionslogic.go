package event

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListEventSubscriptionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListEventSubscriptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListEventSubscriptionsLogic {
	return &ListEventSubscriptionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListEventSubscriptionsLogic) ListEventSubscriptions(req *types.ListEventSubscriptionsReq) (resp *types.ListEventSubscriptionsRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := l.svcCtx.DB.Model(&open_models.OpenAppEventSubscription{}).Where("app_id = ?", req.AppID)
	if req.EventType != "" {
		query = query.Where("event_type = ?", req.EventType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询订阅失败")
	}

	var subs []open_models.OpenAppEventSubscription
	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&subs).Error; err != nil {
		return nil, errors.New("查询订阅失败")
	}

	list := make([]types.ListEventSubscriptionsResItem, 0, len(subs))
	for _, sub := range subs {
		status := sub.Status
		if sub.VerifyStatus != 1 {
			status = 0
		}
		lastVerifiedAt := int64(0)
		if sub.LastVerifiedAt != nil {
			lastVerifiedAt = sub.LastVerifiedAt.Unix()
		}
		list = append(list, types.ListEventSubscriptionsResItem{
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
		})
	}

	return &types.ListEventSubscriptionsRes{
		Total: total,
		List:  list,
	}, nil
}
