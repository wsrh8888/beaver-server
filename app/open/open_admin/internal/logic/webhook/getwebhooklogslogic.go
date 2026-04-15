package webhook

import (
	"context"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWebhookLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Webhook 日志
func NewGetWebhookLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebhookLogsLogic {
	return &GetWebhookLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWebhookLogsLogic) GetWebhookLogs(req *types.GetWebhookLogsReq) (resp *types.GetWebhookLogsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
