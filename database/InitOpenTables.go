package database

import (
	"beaver/app/open/open_models"
	"fmt"

	"gorm.io/gorm"
)

// InitOpenTables 初始化开放平台相关表
func InitOpenTables(db *gorm.DB) error {
	fmt.Println("正在创建开放平台数据表...")

	tables := []interface{}{
		&open_models.OpenApp{},
		&open_models.OpenDeveloper{},
		&open_models.OpenOAuthToken{},
		&open_models.OpenOAuthCode{},
		&open_models.OpenOAuthQrCode{},
		&open_models.OpenEventSubscription{},
		&open_models.OpenBotModel{},
		&open_models.OpenGroupBotModel{},
	}

	for _, table := range tables {
		if err := db.AutoMigrate(table); err != nil {
			return fmt.Errorf("创建表失败: %v", err)
		}
	}

	fmt.Println("开放平台数据表创建成功")
	return nil
}
