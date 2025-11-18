package logic

import (
	"beaver/app/update/update_models"
	"context"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateCityStrategyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新城市策略
func NewUpdateCityStrategyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCityStrategyLogic {
	return &UpdateCityStrategyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCityStrategyLogic) UpdateCityStrategy(req *types.UpdateCityStrategyReq) (resp *types.UpdateCityStrategyRes, err error) {
	// 验证请求参数
	if req.AppID == "" {
		return nil, fmt.Errorf("应用ID不能为空")
	}
	if len(req.CityIDs) == 0 {
		return nil, fmt.Errorf("城市ID列表不能为空")
	}
	if len(req.Strategy) == 0 {
		return nil, fmt.Errorf("策略配置不能为空")
	}

	// 批量更新每个城市的策略
	for _, cityID := range req.CityIDs {
		if err := l.updateSingleCityStrategy(req.AppID, cityID, req.Strategy, req.UpdateType); err != nil {
			logx.Errorf("Failed to update city strategy for city %s: %v", cityID, err)
			return nil, fmt.Errorf("更新城市 %s 策略失败: %v", cityID, err)
		}
	}

	logx.Infof("Successfully updated city strategies for app %s, cities: %v", req.AppID, req.CityIDs)
	return &types.UpdateCityStrategyRes{}, nil
}

// updateSingleCityStrategy 更新单个城市的策略
func (l *UpdateCityStrategyLogic) updateSingleCityStrategy(appID, cityID string, newStrategy []types.UpdateCityStrategyItem, updateType string) error {
	// 查找或创建城市策略记录
	var strategy update_models.UpdateStrategy
	err := l.svcCtx.DB.Where("app_id = ? AND city_id = ?", appID, cityID).First(&strategy).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新的城市策略记录
			strategy = update_models.UpdateStrategy{
				AppID:    appID,
				CityID:   cityID,
				Strategy: &update_models.Strategy{},
				IsActive: true,
			}
		} else {
			return fmt.Errorf("查询城市策略失败: %v", err)
		}
	}

	// 转换新的策略数据
	var newStrategyInfos []update_models.StrategyInfo
	for _, s := range newStrategy {
		newStrategyInfos = append(newStrategyInfos, update_models.StrategyInfo{
			ArchitectureID: s.ArchitectureID,
			VersionID:      s.VersionID,
			ForceUpdate:    s.ForceUpdate,
			IsActive:       s.IsActive,
		})
	}

	var updatedStrategy update_models.Strategy

	if updateType == "global" {
		// 全局更新：直接替换所有策略
		updatedStrategy = update_models.Strategy(newStrategyInfos)
	} else {
		// 单个架构更新：合并现有策略和新策略
		existingStrategy := strategy.Strategy
		if existingStrategy == nil {
			existingStrategy = &update_models.Strategy{}
		}

		// 创建架构ID到策略的映射
		strategyMap := make(map[uint]update_models.StrategyInfo)
		for _, s := range *existingStrategy {
			strategyMap[s.ArchitectureID] = s
		}

		// 更新或添加新策略
		for _, newS := range newStrategyInfos {
			strategyMap[newS.ArchitectureID] = newS
		}

		// 转换回数组
		var mergedStrategy []update_models.StrategyInfo
		for _, s := range strategyMap {
			mergedStrategy = append(mergedStrategy, s)
		}
		updatedStrategy = update_models.Strategy(mergedStrategy)
	}

	// 更新数据库
	if strategy.Id == 0 {
		// 创建新记录
		strategy.Strategy = &updatedStrategy
		if err := l.svcCtx.DB.Create(&strategy).Error; err != nil {
			return fmt.Errorf("创建城市策略失败: %v", err)
		}
	} else {
		// 更新现有记录
		updateData := map[string]interface{}{
			"strategy": &updatedStrategy,
		}
		if err := l.svcCtx.DB.Model(&strategy).Updates(updateData).Error; err != nil {
			return fmt.Errorf("更新城市策略失败: %v", err)
		}
	}

	logx.Infof("Successfully updated city strategy for app %s, city %s, updateType: %s", appID, cityID, updateType)
	return nil
}
