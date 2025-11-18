package logic

import (
	"beaver/app/update/update_models"
	"context"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCityStrategiesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取城市策略列表
func NewGetCityStrategiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCityStrategiesLogic {
	return &GetCityStrategiesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCityStrategiesLogic) GetCityStrategies(req *types.GetCityStrategiesReq) (resp *types.GetCityStrategiesRes, err error) {
	// 构建查询条件
	query := l.svcCtx.DB.Model(&update_models.UpdateStrategy{})

	// 应用ID过滤
	if req.AppID != "" {
		query = query.Where("app_id = ?", req.AppID)
	}

	// 活跃状态过滤
	if req.IsActive {
		query = query.Where("is_active = ?", true)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("Failed to count city strategies: %v", err)
		return nil, fmt.Errorf("获取城市策略总数失败")
	}

	// 分页
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize

	// 查询策略列表
	var strategies []update_models.UpdateStrategy
	if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&strategies).Error; err != nil {
		logx.Errorf("Failed to get city strategies: %v", err)
		return nil, fmt.Errorf("获取城市策略列表失败")
	}

	// 转换为响应格式
	var strategyInfos []types.GetCityStrategiesItem
	for _, strategy := range strategies {
		// 转换策略信息
		var strategyInfosList []types.GetCityStrategiesStrategyItem
		if strategy.Strategy != nil {
			for _, s := range *strategy.Strategy {
				// 查询版本信息
				var version update_models.UpdateVersion
				versionStr := "未知版本"
				if err := l.svcCtx.DB.Where("id = ?", s.VersionID).First(&version).Error; err == nil {
					versionStr = version.Version
				}

				strategyInfosList = append(strategyInfosList, types.GetCityStrategiesStrategyItem{
					ArchitectureID: s.ArchitectureID,
					VersionID:      s.VersionID,
					Version:        versionStr,
					ForceUpdate:    s.ForceUpdate,
					IsActive:       s.IsActive,
				})
			}
		}

		strategyInfo := types.GetCityStrategiesItem{
			Id:        uint(strategy.Id),
			AppID:     strategy.AppID,
			CityID:    strategy.CityID,
			Strategy:  strategyInfosList,
			IsActive:  strategy.IsActive,
			CreatedAt: strategy.CreatedAt.String(),
			UpdatedAt: strategy.UpdatedAt.String(),
		}
		strategyInfos = append(strategyInfos, strategyInfo)
	}

	return &types.GetCityStrategiesRes{
		Total:      total,
		Strategies: strategyInfos,
	}, nil
}
