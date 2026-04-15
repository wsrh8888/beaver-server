package stats

import (
	"context"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAPICallsStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 API 调用统计
func NewGetAPICallsStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAPICallsStatsLogic {
	return &GetAPICallsStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAPICallsStatsLogic) GetAPICallsStats(req *types.GetAPICallsStatsReq) (resp *types.GetAPICallsStatsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
