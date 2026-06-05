package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCityStrategyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCityStrategyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCityStrategyLogic {
	return &UpdateCityStrategyLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateCityStrategyLogic) UpdateCityStrategy(req *types.UpdateCityStrategyReq) (resp *types.UpdateCityStrategyRes, err error) {
	inputs := make([]*platform_rpc.StrategyInput, 0, len(req.Strategy))
	for _, s := range req.Strategy {
		inputs = append(inputs, &platform_rpc.StrategyInput{
			ArchitectureId: uint64(s.ArchitectureID),
			VersionId:      uint64(s.VersionID),
			ForceUpdate:    s.ForceUpdate,
			IsActive:       s.IsActive,
		})
	}

	_, err = l.svcCtx.PlatformRpc.UpdateCityStrategy(l.ctx, &platform_rpc.UpdateCityStrategyReq{
		AppId:      req.AppID,
		CityIds:    req.CityIDs,
		Strategy:   inputs,
		UpdateType: req.UpdateType,
	})
	if err != nil {
		l.Errorf("更新城市策略失败: %v", err)
		return nil, err
	}
	return &types.UpdateCityStrategyRes{}, nil
}
