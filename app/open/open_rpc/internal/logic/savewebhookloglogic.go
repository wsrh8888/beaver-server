package logic

import (
	"context"

	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type SaveWebhookLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveWebhookLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveWebhookLogLogic {
	return &SaveWebhookLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("SaveWebhookLogLogic", ctx),
	}
}

func (l *SaveWebhookLogLogic) SaveWebhookLog(in *open_rpc.SaveWebhookLogReq) (*open_rpc.SaveWebhookLogRes, error) {
	// todo: add your logic here and delete this line

	return &open_rpc.SaveWebhookLogRes{}, nil
}
