package permission

import (
	"context"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigAppPermissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfigAppPermissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigAppPermissionLogic {
	return &ConfigAppPermissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigAppPermissionLogic) ConfigAppPermission(req *types.ConfigAppPermissionReq) (resp *types.ConfigAppPermissionRes, err error) {
	// TODO: 实现配置应用权限逻辑
	return &types.ConfigAppPermissionRes{}, nil
}
