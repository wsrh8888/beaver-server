package event

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWebhookListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWebhookListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebhookListLogic {
	return &GetWebhookListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetWebhookListLogic) GetWebhookList(req *types.GetWebhookListReq, authorization string) (resp *types.GetWebhookListRes, err error) {
	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := l.svcCtx.DB.Model(&open_models.OpenAppEventSubscription{}).Where("app_id = ?", token.AppID)
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

	list := make([]types.WebhookSubscriptionInfo, 0, len(subs))
	for _, sub := range subs {
		status := sub.Status
		if sub.VerifyStatus != 1 {
			status = 0
		}
		lastVerifiedAt := int64(0)
		if sub.LastVerifiedAt != nil {
			lastVerifiedAt = sub.LastVerifiedAt.Unix()
		}
		list = append(list, types.WebhookSubscriptionInfo{
			ID:             uint64(sub.ID),
			EventType:      sub.EventType,
			TargetURL:      sub.CallbackURL,
			Status:         status,
			VerifyStatus:   sub.VerifyStatus,
			LastError:      sub.LastError,
			LastVerifiedAt: lastVerifiedAt,
			RetryCount:     sub.RetryCount,
			Timeout:        sub.Timeout,
			CreatedAt:      sub.CreatedAt.Unix(),
		})
	}

	return &types.GetWebhookListRes{Total: total, List: list}, nil
}
