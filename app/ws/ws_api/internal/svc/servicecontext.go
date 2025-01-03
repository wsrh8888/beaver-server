package svc

import (
	"beaver/app/ws/ws_api/internal/config"
	"beaver/core"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Client
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)

	return &ServiceContext{
		DB:     mysqlDb,
		Redis:  client,
		Config: c,
	}
}
