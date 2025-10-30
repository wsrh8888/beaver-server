package database

import (
	"beaver/app/update/update_models"
	"beaver/utils/conversation"
	"fmt"

	"gorm.io/gorm"
)

// 初始化城市策略数据
func InitUpdateData(db *gorm.DB) error {
	fmt.Println("=== 开始初始化城市策略数据 ===")

	// 从数据库读取所有应用，为每个应用创建城市策略
	if err := initCityStrategyData(db); err != nil {
		return fmt.Errorf("初始化城市策略数据失败: %v", err)
	}

	fmt.Println("=== 城市策略数据初始化完成 ===")
	return nil
}

// 初始化城市策略数据
func initCityStrategyData(db *gorm.DB) error {
	fmt.Println("正在初始化城市策略数据...")

	// 从数据库读取所有应用
	var apps []update_models.UpdateApp
	if err := db.Where("is_active = ?", true).Find(&apps).Error; err != nil {
		return fmt.Errorf("读取应用数据失败: %v", err)
	}

	if len(apps) == 0 {
		fmt.Println("警告: 数据库中没有找到活跃的应用")
		return nil
	}

	cities := conversation.GetDefaultCities()

	// 为每个应用创建所有城市的默认策略
	for _, app := range apps {
		for _, city := range cities {
			var count int64
			// 检查城市策略是否存在
			db.Model(&update_models.UpdateStrategy{}).
				Where("app_id = ? AND city_id = ?", app.UUID, city.Code).
				Count(&count)

			if count == 0 {
				// 城市策略不存在，创建默认策略
				defaultStrategy := &update_models.Strategy{}

				newStrategy := update_models.UpdateStrategy{
					AppID:    app.UUID,
					CityID:   city.Code,
					Strategy: defaultStrategy,
					IsActive: true,
				}

				if err := db.Create(&newStrategy).Error; err != nil {
					return fmt.Errorf("创建城市策略失败 (App: %s, City: %s): %v", app.UUID, city.Code, err)
				}
				fmt.Printf("已创建城市策略: %s - %s (%s)\n", app.Name, city.Name, city.Code)
			} else {
				fmt.Printf("城市策略已存在，跳过: %s - %s (%s)\n", app.Name, city.Name, city.Code)
			}
		}
	}

	fmt.Printf("成功处理 %d 个应用，%d 个城市的策略数据\n", len(apps), len(cities))
	return nil
}
