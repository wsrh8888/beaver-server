package logic

import (
	"context"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListCityStrategiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCityStrategiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCityStrategiesLogic {
	return &ListCityStrategiesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListCityStrategiesLogic) ListCityStrategies(in *platform_rpc.ListCityStrategiesReq) (*platform_rpc.ListCityStrategiesRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&platform_models.UpdateStrategy{})
	if in.AppId != "" {
		db = db.Where("app_id = ?", in.AppId)
	}
	if in.IsActive != nil && *in.IsActive {
		db = db.Where("is_active = ?", true)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计城市策略失败: %v", err)
		return nil, status.Error(codes.Internal, "获取城市策略总数失败")
	}

	var list []platform_models.UpdateStrategy
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&list).Error; err != nil {
		l.Errorf("查询城市策略失败: %v", err)
		return nil, status.Error(codes.Internal, "获取城市策略列表失败")
	}

	items := make([]*platform_rpc.CityStrategyItem, 0, len(list))
	for _, s := range list {
		strategyItems := make([]*platform_rpc.StrategyInfoItem, 0)
		if s.Strategy != nil {
			for _, info := range *s.Strategy {
				versionStr := "未知版本"
				var version platform_models.UpdateVersion
				if err := l.svcCtx.DB.Where("id = ?", info.VersionID).First(&version).Error; err == nil {
					versionStr = version.Version
				}
				strategyItems = append(strategyItems, &platform_rpc.StrategyInfoItem{
					ArchitectureId: uint64(info.ArchitectureID),
					VersionId:      uint64(info.VersionID),
					Version:        versionStr,
					ForceUpdate:    info.ForceUpdate,
					IsActive:       info.IsActive,
				})
			}
		}

		items = append(items, &platform_rpc.CityStrategyItem{
			Id:        uint64(s.Id),
			AppId:     s.AppID,
			CityId:    s.CityID,
			Strategy:  strategyItems,
			IsActive:  s.IsActive,
			CreatedAt: s.CreatedAt.String(),
			UpdatedAt: s.UpdatedAt.String(),
		})
	}

	return &platform_rpc.ListCityStrategiesRes{Total: total, Strategies: items}, nil
}
