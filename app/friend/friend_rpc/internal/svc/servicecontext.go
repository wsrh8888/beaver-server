package svc

import (
	"beaver/app/friend/friend_rpc/internal/config"
	"beaver/core"

	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)

	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
