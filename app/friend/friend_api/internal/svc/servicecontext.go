package svc

import (
	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/friend/friend_api/internal/config"
	"beaver/app/friend/friend_rpc/friend"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/notification/notification_rpc/notification"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"
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
	NotifyRpc  notification.Notification
	VersionGen *versionPkg.VersionGenerator
	RocketMQ   *corerocketmq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	// 初始化 RocketMQ 客户端

	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)

	return &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		Redis:      client,
		ChatRpc:    chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FriendRpc:  friend.NewFriend(zrpc.MustNewClient(c.FriendRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		NotifyRpc:  notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		VersionGen: versionGen,
		RocketMQ:   mqClient,
	}
}
