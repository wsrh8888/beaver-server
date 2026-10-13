package logic

import (
	"context"

	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type ListWebhookLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListWebhookLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWebhookLogsLogic {
	return &ListWebhookLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListWebhookLogsLogic", ctx),
	}
}

func (l *ListWebhookLogsLogic) ListWebhookLogs(in *open_rpc.ListWebhookLogsReq) (*open_rpc.ListWebhookLogsRes, error) {
	// todo: add your logic here and delete this line

	return &open_rpc.ListWebhookLogsRes{}, nil
}
