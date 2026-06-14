package logic

import (
	"context"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveWebhookLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveWebhookLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveWebhookLogLogic {
	return &SaveWebhookLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SaveWebhookLogLogic) SaveWebhookLog(in *open_rpc.SaveWebhookLogReq) (*open_rpc.SaveWebhookLogRes, error) {
	log := open_models.OpenWebhookLog{
		ConfigID:  in.ConfigId,
		AppID:     in.AppId,
		EventType: in.EventType,
		Status:    int(in.Status),
	}
	if err := l.svcCtx.DB.Create(&log).Error; err != nil {
		l.Errorf("保存 Webhook 日志失败: %v", err)
		return nil, err
	}
	return &open_rpc.SaveWebhookLogRes{}, nil
}
