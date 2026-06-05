package svc

import (
	"beaver/app/platform/platform_rpc/internal/config"
	"beaver/core/coregorm"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		DB:     coregorm.InitGorm(c.Mysql.DataSource),
	}
}
