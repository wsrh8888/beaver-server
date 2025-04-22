package svc

import (
	"beaver/app/group/group_rpc/internal/config"
	"beaver/core"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)

	return &ServiceContext{
		Config: c,
		Redis:  client,
		DB:     mysqlDb,
	}
}
