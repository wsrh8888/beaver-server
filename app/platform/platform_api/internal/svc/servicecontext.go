package svc

import (
	"beaver/app/file/file_rpc/file"
	"beaver/app/file/file_rpc/types/file_rpc"
	platformcli "beaver/app/platform/platform_rpc/platform"
	"beaver/app/platform/platform_api/internal/config"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	FileRpc     file_rpc.FileClient
	PlatformRpc platformcli.Platform
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	return &ServiceContext{
		Config:      c,
		DB:          mysqlDb,
		Redis:       client,
		FileRpc:     file.NewFile(zrpc.MustNewClient(c.FileRpc, rpcOpt)),
		PlatformRpc: platformcli.NewPlatform(zrpc.MustNewClient(c.PlatformRpc, rpcOpt)),
	}
}
