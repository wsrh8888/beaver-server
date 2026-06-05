package logic

import (
	"context"
	"errors"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UpdateCityStrategyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCityStrategyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCityStrategyLogic {
	return &UpdateCityStrategyLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateCityStrategyLogic) UpdateCityStrategy(in *platform_rpc.UpdateCityStrategyReq) (*platform_rpc.UpdateCityStrategyRes, error) {
	if in.AppId == "" {
		return nil, status.Error(codes.InvalidArgument, "应用ID不能为空")
	}
	if len(in.CityIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "城市ID列表不能为空")
	}
	if len(in.Strategy) == 0 {
		return nil, status.Error(codes.InvalidArgument, "策略配置不能为空")
	}

	for _, cityID := range in.CityIds {
		if err := l.updateSingle(in.AppId, cityID, in.Strategy, in.UpdateType); err != nil {
			l.Errorf("更新城市策略失败 city=%s: %v", cityID, err)
			return nil, status.Errorf(codes.Internal, "更新城市 %s 策略失败", cityID)
		}
	}

	return &platform_rpc.UpdateCityStrategyRes{}, nil
}

func (l *UpdateCityStrategyLogic) updateSingle(appID, cityID string, inputs []*platform_rpc.StrategyInput, updateType string) error {
	var strategy platform_models.UpdateStrategy
	err := l.svcCtx.DB.Where("app_id = ? AND city_id = ?", appID, cityID).First(&strategy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			strategy = platform_models.UpdateStrategy{
				AppID:    appID,
				CityID:   cityID,
				Strategy: &platform_models.Strategy{},
				IsActive: true,
			}
		} else {
			return err
		}
	}

	newInfos := make([]platform_models.StrategyInfo, 0, len(inputs))
	for _, s := range inputs {
		newInfos = append(newInfos, platform_models.StrategyInfo{
			ArchitectureID: uint(s.ArchitectureId),
			VersionID:      uint(s.VersionId),
			ForceUpdate:    s.ForceUpdate,
			IsActive:       s.IsActive,
		})
	}

	var updated platform_models.Strategy
	if updateType == "global" {
		updated = platform_models.Strategy(newInfos)
	} else {
		existing := strategy.Strategy
		if existing == nil {
			existing = &platform_models.Strategy{}
		}
		strategyMap := make(map[uint]platform_models.StrategyInfo)
		for _, s := range *existing {
			strategyMap[s.ArchitectureID] = s
		}
		for _, s := range newInfos {
			strategyMap[s.ArchitectureID] = s
		}
		merged := make([]platform_models.StrategyInfo, 0, len(strategyMap))
		for _, s := range strategyMap {
			merged = append(merged, s)
		}
		updated = platform_models.Strategy(merged)
	}

	if strategy.Id == 0 {
		strategy.Strategy = &updated
		return l.svcCtx.DB.Create(&strategy).Error
	}

	return l.svcCtx.DB.Model(&strategy).Updates(map[string]interface{}{"strategy": &updated}).Error
}
