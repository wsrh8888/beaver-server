package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCityStrategiesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCityStrategiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCityStrategiesLogic {
	return &GetCityStrategiesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetCityStrategiesLogic) GetCityStrategies(req *types.GetCityStrategiesReq) (resp *types.GetCityStrategiesRes, err error) {
	rpcReq := &platform_rpc.ListCityStrategiesReq{
		AppId:    req.AppID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	if req.IsActive {
		active := true
		rpcReq.IsActive = &active
	}

	rpcRes, err := l.svcCtx.PlatformRpc.ListCityStrategies(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("获取城市策略失败: %v", err)
		return nil, err
	}

	list := make([]types.GetCityStrategiesItem, 0, len(rpcRes.Strategies))
	for _, s := range rpcRes.Strategies {
		strategyList := make([]types.GetCityStrategiesStrategyItem, 0, len(s.Strategy))
		for _, info := range s.Strategy {
			strategyList = append(strategyList, types.GetCityStrategiesStrategyItem{
				ArchitectureID: uint(info.ArchitectureId),
				VersionID:      uint(info.VersionId),
				Version:        info.Version,
				ForceUpdate:    info.ForceUpdate,
				IsActive:       info.IsActive,
			})
		}
		list = append(list, types.GetCityStrategiesItem{
			Id:        uint(s.Id),
			AppID:     s.AppId,
			CityID:    s.CityId,
			Strategy:  strategyList,
			IsActive:  s.IsActive,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		})
	}

	return &types.GetCityStrategiesRes{Total: rpcRes.Total, Strategies: list}, nil
}
