package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateArchitectureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArchitectureLogic {
	return &UpdateArchitectureLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateArchitectureLogic) UpdateArchitecture(req *types.UpdateArchitectureReq) (resp *types.UpdateArchitectureRes, err error) {
	_, err = l.svcCtx.PlatformRpc.UpdateArchitecture(l.ctx, &platform_rpc.UpdateArchitectureReq{
		Id:          uint64(req.Id),
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		l.Errorf("更新架构失败: %v", err)
		return nil, err
	}
	return &types.UpdateArchitectureRes{}, nil
}
