// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package webhook

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebhookSendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 通过 Webhook 发送消息（无需鉴权，通过 URL 中的 token 验证）
func NewWebhookSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebhookSendMessageLogic {
	return &WebhookSendMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebhookSendMessageLogic) WebhookSendMessage(req *types.WebhookSendMessageReq) (resp *types.WebhookSendMessageRes, err error) {
	// todo: add your logic here and delete this line

	return
}
