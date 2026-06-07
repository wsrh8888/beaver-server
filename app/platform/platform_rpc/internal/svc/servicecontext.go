package svc

import (
	"beaver/app/platform/platform_rpc/internal/config"
	"beaver/core/coregorm"
	"beaver/database"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	_ = database.SeedUpdateData(db)
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
