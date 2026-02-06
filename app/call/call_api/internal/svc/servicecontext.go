package svc

import (
	"beaver/app/call/call_api/internal/config"
	"beaver/app/call/call_rpc/call"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
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
	UserRpc user_rpc.UserClient
	CallRpc call_rpc.CallClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)

	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   client,
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		CallRpc: call.NewCall(zrpc.MustNewClient(c.CallRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
