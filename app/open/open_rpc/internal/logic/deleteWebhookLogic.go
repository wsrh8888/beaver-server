package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWebhookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWebhookLogic {
	return &DeleteWebhookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteWebhookLogic) DeleteWebhook(in *open_rpc.DeleteWebhookReq) (*open_rpc.DeleteWebhookRes, error) {
	var record open_models.OpenGroupBotModel
	if err := l.svcCtx.DB.First(&record, in.Id).Error; err != nil {
		return nil, errors.New("webhook 不存在")
	}
	if err := l.svcCtx.DB.Delete(&record).Error; err != nil {
		return nil, errors.New("删除失败")
	}
	return &open_rpc.DeleteWebhookRes{}, nil
}
