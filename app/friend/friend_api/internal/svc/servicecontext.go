package svc

import (
	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/friend/friend_api/internal/config"
	"beaver/app/friend/friend_rpc/friend"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core"
	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *redis.Client
	UserRpc    user_rpc.UserClient
	ChatRpc    chat_rpc.ChatClient
	FriendRpc  friend_rpc.FriendClient
	VersionGen *versionPkg.VersionGenerator
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	return &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		Redis:      client,
		ChatRpc:    chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FriendRpc:  friend.NewFriend(zrpc.MustNewClient(c.FriendRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		VersionGen: versionGen,
	}
}
