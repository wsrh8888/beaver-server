package database

import (
	"beaver/app/open/open_models"

	"gorm.io/gorm"
)

// InitOpenTables 初始化开放平台相关表
func InitOpenTables(db *gorm.DB) {
	db.AutoMigrate(
		&open_models.OpenApp{},
		&open_models.OpenAppVersion{},
		&open_models.OpenAppPermission{},
		&open_models.OpenAccessToken{},
		&open_models.OpenRefreshToken{},
		&open_models.OpenWebhookConfig{},
		&open_models.OpenWebhookLog{},
		&open_models.OpenAPILog{},
		&open_models.OpenDeveloper{}, // 开发者申请表
	)
}
