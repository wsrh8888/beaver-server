package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateWebhookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateWebhookLogic {
	return &UpdateWebhookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateWebhookLogic) UpdateWebhook(in *open_rpc.UpdateWebhookReq) (*open_rpc.UpdateWebhookRes, error) {
	var record open_models.OpenGroupBotModel
	if err := l.svcCtx.DB.First(&record, in.Id).Error; err != nil {
		return nil, errors.New("bot 不存在")
	}

	if in.Status == -1 {
		return &open_rpc.UpdateWebhookRes{}, nil
	}
	if err := l.svcCtx.DB.Model(&record).Update("status", in.Status).Error; err != nil {
		return nil, errors.New("更新失败")
	}
	return &open_rpc.UpdateWebhookRes{}, nil
}
