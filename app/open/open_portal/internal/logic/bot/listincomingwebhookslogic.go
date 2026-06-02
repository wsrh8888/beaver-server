package bot

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListIncomingWebhooksLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Incoming Webhook 列表
func NewListIncomingWebhooksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListIncomingWebhooksLogic {
	return &ListIncomingWebhooksLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListIncomingWebhooksLogic) ListIncomingWebhooks(req *types.ListIncomingWebhooksReq) (resp *types.ListIncomingWebhooksRes, err error) {
	if _, err := l.svcCtx.RequireDeveloper(req.UserID); err != nil {
		return nil, err
	}

	// todo: add your logic here and delete this line

	return
}
