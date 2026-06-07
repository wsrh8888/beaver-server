package monitor

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在线用户统计
func NewGetOnlineStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineStatsLogic {
	return &GetOnlineStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOnlineStatsLogic) GetOnlineStats(req *types.GetOnlineStatsReq) (resp *types.GetOnlineStatsRes, err error) {
	// todo: add your logic here and delete this line

	return
}
