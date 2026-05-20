package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 注册 Bot Webhook URL（注册后 Beaver 立即向该 URL 发送 Challenge 验证请求）
func NewRegisterWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterWebhookLogic {
	return &RegisterWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterWebhookLogic) RegisterWebhook(req *types.RegisterWebhookReq) (resp *types.RegisterWebhookRes, err error) {
	// todo: add your logic here and delete this line

	return
}
