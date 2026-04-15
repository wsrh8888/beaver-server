package webhook

import (
	"context"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 配置 Webhook
func NewConfigWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigWebhookLogic {
	return &ConfigWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigWebhookLogic) ConfigWebhook(req *types.ConfigWebhookReq) (resp *types.ConfigWebhookRes, err error) {
	// todo: add your logic here and delete this line

	return
}
