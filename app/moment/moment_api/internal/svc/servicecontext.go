package svc

import (
	"beaver/app/friend/friend_rpc/friend"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/moment/moment_api/internal/config"
	"beaver/app/notification/notification_rpc/notification"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
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

	UserRpc   user_rpc.UserClient
	NotifyRpc notification.Notification
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:    c,
		DB:        mysqlDb,
		Redis:     client,
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FriendRpc: friend.NewFriend(zrpc.MustNewClient(c.FriendRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		NotifyRpc: notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
