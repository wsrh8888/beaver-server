package svc

import (
	"beaver/app/call/call_api/internal/config"
	"beaver/app/call/call_rpc/call"
	"beaver/app/call/call_rpc/types/call_rpc"
	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/group"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	UserRpc  user_rpc.UserClient
	CallRpc  call_rpc.CallClient
	ChatRpc  chat_rpc.ChatClient
	GroupRpc group_rpc.GroupClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	return &ServiceContext{
		Config:   c,
		DB:       mysqlDb,
		Redis:    rdb,
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		CallRpc:  call.NewCall(zrpc.MustNewClient(c.CallRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		ChatRpc:  chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		GroupRpc: group.NewGroup(zrpc.MustNewClient(c.GroupRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
