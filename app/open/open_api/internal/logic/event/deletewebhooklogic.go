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

type DeleteWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWebhookLogic {
	return &DeleteWebhookLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteWebhookLogic) DeleteWebhook(req *types.DeleteWebhookReq, authorization string) (resp *types.DeleteWebhookRes, err error) {
	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}

	result := l.svcCtx.DB.Where("id = ? AND app_id = ?", req.SubscriptionID, token.AppID).
		Delete(&open_models.OpenAppEventSubscription{})
	if result.Error != nil {
		return nil, errors.New("删除订阅失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("订阅不存在")
	}

	return &types.DeleteWebhookRes{Success: true}, nil
}
