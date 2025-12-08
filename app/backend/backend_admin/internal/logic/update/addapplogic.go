package logic

import (
	"beaver/app/dictionary/dictionary_rpc/types/dictionary_rpc"
	"beaver/app/update/update_models"
	"context"
	"errors"
	"fmt"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加新应用
func NewAddAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppLogic {
	return &AddAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAppLogic) AddApp(req *types.AddAppReq) (resp *types.AddAppRes, err error) {
	// 检查应用名称是否已存在
	var existingApp update_models.UpdateApp
	if err := l.svcCtx.DB.Where("name = ?", req.Name).First(&existingApp).Error; err == nil {
		return nil, fmt.Errorf("应用名称已存在")
	}

	// 创建新应用
	app := update_models.UpdateApp{
		Name:        req.Name,
		AppID:       strings.Replace(uuid.New().String(), "-", "", -1),
		Description: req.Description,
		IsActive:    true, // 默认为活跃状态
		Icon:        "",
	}

	if err := l.svcCtx.DB.Create(&app).Error; err != nil {
		logx.Errorf("Failed to create app: %v", err)
		return nil, errors.New("创建应用失败")
	}

	// 为新应用自动初始化所有城市的策略
	if err := l.initCityStrategiesForApp(app.AppID); err != nil {
		logx.Errorf("Failed to init city strategies for app %s: %v", app.AppID, err)
		// 不返回错误，因为应用创建成功了，只是策略初始化失败
	}

	return &types.AddAppRes{
		Id:    uint(app.Id),
		AppID: app.AppID,
	}, nil
}

// 为新应用初始化所有城市的策略
func (l *AddAppLogic) initCityStrategiesForApp(appID string) error {
	// 通过 RPC 调用获取城市列表
	citiesRes, err := l.svcCtx.DictionaryRpc.GetCities(l.ctx, &dictionary_rpc.GetCitiesReq{})
	if err != nil {
		return fmt.Errorf("获取城市列表失败: %v", err)
	}

	for _, city := range citiesRes.Cities {
		// 检查城市策略是否已存在
		var count int64
		l.svcCtx.DB.Model(&update_models.UpdateStrategy{}).
			Where("app_id = ? AND city_id = ?", appID, city.CityId).
			Count(&count)

		if count == 0 {
			// 城市策略不存在，创建默认策略
			defaultStrategy := &update_models.Strategy{}

			newStrategy := update_models.UpdateStrategy{
				AppID:    appID,
				CityID:   city.CityId,
				Strategy: defaultStrategy,
				IsActive: true,
			}

			if err := l.svcCtx.DB.Create(&newStrategy).Error; err != nil {
				return fmt.Errorf("创建城市策略失败 (App: %s, City: %s): %v", appID, city.CityId, err)
			}
			logx.Infof("已为应用 %s 创建城市策略: %s (%s)", appID, city.CityName, city.CityId)
		}
	}

	logx.Infof("成功为应用 %s 初始化了 %d 个城市的策略", appID, len(citiesRes.Cities))
	return nil
}
