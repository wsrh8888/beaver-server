package svc

import (
	"beaver/app/file/file_rpc/file"
	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/app/update/update_api/internal/config"
	"beaver/common/zrpc_interceptor"
	"beaver/core"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redis.Client
	FileRpc file_rpc.FileClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   client,
		FileRpc: file.NewFile(zrpc.MustNewClient(c.FileRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
