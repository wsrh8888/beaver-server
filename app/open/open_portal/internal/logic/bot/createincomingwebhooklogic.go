package bot

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateIncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建 Incoming Webhook（开放平台应用维度）
func NewCreateIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateIncomingWebhookLogic {
	return &CreateIncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateIncomingWebhookLogic) CreateIncomingWebhook(req *types.CreateIncomingWebhookReq) (resp *types.CreateIncomingWebhookRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// TODO: OpenGroupBotModel 已移除，此功能暂时禁用
	return nil, errors.New("Incoming Webhook 功能暂未实现")
}
