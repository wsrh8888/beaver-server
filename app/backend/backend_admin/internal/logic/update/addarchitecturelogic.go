package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddArchitectureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddArchitectureLogic {
	return &AddArchitectureLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddArchitectureLogic) AddArchitecture(req *types.AddArchitectureReq) (resp *types.AddArchitectureRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateArchitecture(l.ctx, &platform_rpc.CreateArchitectureReq{
		AppId:       req.AppID,
		PlatformId:  uint32(req.PlatformID),
		ArchId:      uint32(req.ArchID),
		Description: req.Description,
	})
	if err != nil {
		l.Errorf("创建架构失败: %v", err)
		return nil, err
	}
	return &types.AddArchitectureRes{Id: uint(rpcRes.Id)}, nil
}
