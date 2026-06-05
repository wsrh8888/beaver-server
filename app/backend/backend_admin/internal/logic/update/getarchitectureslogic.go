package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetArchitecturesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArchitecturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArchitecturesLogic {
	return &GetArchitecturesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetArchitecturesLogic) GetArchitectures(req *types.GetArchitecturesReq) (resp *types.GetArchitecturesRes, err error) {
	rpcReq := &platform_rpc.ListArchitecturesReq{
		AppId:    req.AppID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	if req.IsActive {
		active := true
		rpcReq.IsActive = &active
	}

	rpcRes, err := l.svcCtx.PlatformRpc.ListArchitectures(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("获取架构列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetArchitecturesItem, 0, len(rpcRes.Architectures))
	for _, arch := range rpcRes.Architectures {
		list = append(list, types.GetArchitecturesItem{
			Id:          uint(arch.Id),
			AppID:       arch.AppId,
			AppName:     arch.AppName,
			PlatformID:  uint(arch.PlatformId),
			ArchID:      uint(arch.ArchId),
			Description: arch.Description,
			IsActive:    arch.IsActive,
			CreatedAt:   arch.CreatedAt,
			UpdatedAt:   arch.UpdatedAt,
		})
	}

	return &types.GetArchitecturesRes{Total: rpcRes.Total, Architectures: list}, nil
}
