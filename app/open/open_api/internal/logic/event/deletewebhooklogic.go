package event

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除 Webhook 订阅
func NewDeleteWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWebhookLogic {
	return &DeleteWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteWebhookLogic) DeleteWebhook(req *types.DeleteWebhookReq) (resp *types.DeleteWebhookRes, err error) {
	// todo: add your logic here and delete this line

	return
}
