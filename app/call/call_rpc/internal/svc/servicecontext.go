package svc

import (
	"beaver/app/call/call_rpc/internal/config"
	"beaver/core"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisConf.Addr,
		Password: c.RedisConf.Password,
		DB:       c.RedisConf.Db,
	})

	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
		Redis:  rdb,
	}
}
