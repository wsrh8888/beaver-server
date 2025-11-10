package svc

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_api/internal/config"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	Redis     *redis.Client
	FriendRpc friend_rpc.FriendClient
	GroupRpc  group_rpc.GroupClient
	UserRpc   user_rpc.UserClient
	ChatRpc   chat_rpc.ChatClient
	MomentRpc moment_rpc.MomentClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:    c,
		DB:        mysqlDb,
		Redis:     client,
		FriendRpc: friend_rpc.NewFriendClient(zrpc.MustNewClient(c.FriendRpc).Conn()),
		GroupRpc:  group_rpc.NewGroupClient(zrpc.MustNewClient(c.GroupRpc).Conn()),
		UserRpc:   user_rpc.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
		ChatRpc:   chat_rpc.NewChatClient(zrpc.MustNewClient(c.ChatRpc).Conn()),
		MomentRpc: moment_rpc.NewMomentClient(zrpc.MustNewClient(c.MomentRpc).Conn()),
	}
}
