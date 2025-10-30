package svc

import (
	"beaver/app/datasync/datasync_api/internal/config"
	"beaver/app/datasync/datasync_rpc/types/types/datasync_rpc"
	"beaver/core"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	DatasyncRpc datasync_rpc.DatasyncClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:      c,
		DB:          mysqlDb,
		Redis:       client,
		DatasyncRpc: datasync_rpc.NewDatasyncClient(zrpc.MustNewClient(c.DatasyncRpc).Conn()),
	}
}
