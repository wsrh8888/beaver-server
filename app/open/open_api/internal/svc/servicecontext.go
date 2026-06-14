package svc

import (
	"beaver/app/auth/auth_rpc/auth"
	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/group/group_rpc/group"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_api/internal/config"
	oauthmiddle "beaver/app/open/open_api/internal/middle/oauth"
	"beaver/app/open/open_rpc/open"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	OAuth    *oauthmiddle.Qrcode
	UserRpc  user_rpc.UserClient
	AuthRpc  auth.Auth
	ChatRpc  chat_rpc.ChatClient
	GroupRpc group_rpc.GroupClient
	OpenRpc  open_rpc.OpenClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)

	return &ServiceContext{
		Config:   c,
		DB:       mysqlDb,
		Redis:    client,
		OAuth:    oauthmiddle.NewQrcode(mysqlDb),
		UserRpc:  user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		AuthRpc:  auth.NewAuth(zrpc.MustNewClient(c.AuthRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		ChatRpc:  chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		GroupRpc: group.NewGroup(zrpc.MustNewClient(c.GroupRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		OpenRpc:  open.NewOpen(zrpc.MustNewClient(c.OpenRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
