package quota

import (
	"context"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigQuotaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigQuotaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigQuotaLogic {
	return &ConfigQuotaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigQuotaLogic) ConfigQuota(req *types.ConfigQuotaReq) (resp *types.ConfigQuotaRes, err error) {
	// TODO: 实现配置配额逻辑
	return &types.ConfigQuotaRes{
		Success: true,
	}, nil
}
