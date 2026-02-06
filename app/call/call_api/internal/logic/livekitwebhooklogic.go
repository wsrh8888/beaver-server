package logic

import (
	"context"

	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LiveKitWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// LiveKit 服务器回调 (需在网关配置白名单)
func NewLiveKitWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LiveKitWebhookLogic {
	return &LiveKitWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LiveKitWebhookLogic) LiveKitWebhook(req *types.LiveKitWebhookReq) (resp *types.LiveKitWebhookRes, err error) {
	// todo: add your logic here and delete this line

	return
}
