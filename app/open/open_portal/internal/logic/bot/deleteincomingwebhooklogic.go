package bot

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteIncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除 Incoming Webhook
func NewDeleteIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteIncomingWebhookLogic {
	return &DeleteIncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteIncomingWebhookLogic) DeleteIncomingWebhook(req *types.DeleteIncomingWebhookReq) (resp *types.DeleteIncomingWebhookRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// todo: add your logic here and delete this line

	return
}
