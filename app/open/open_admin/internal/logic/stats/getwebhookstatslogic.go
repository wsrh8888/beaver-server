package stats

import (
	"context"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWebhookStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Webhook 统计
func NewGetWebhookStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWebhookStatsLogic {
	return &GetWebhookStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWebhookStatsLogic) GetWebhookStats(req *types.GetWebhookStatsReq) (resp *types.GetWebhookStatsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
