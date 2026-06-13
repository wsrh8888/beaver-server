package svc

import (
	"beaver/app/auth/auth_api/internal/config"
	"beaver/app/notification/notification_rpc/notification"
	"beaver/app/open/open_rpc/open"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config           config.Config
	Redis            *redis.Client
	DB               *gorm.DB
	UserRpc          user_rpc.UserClient
	NotificationRpc  notification.Notification
	OpenRpc          open.Open
	RocketMQ         *corerocketmq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)

	return &ServiceContext{
		Config:          c,
		Redis:           client,
		DB:              mysqlDb,
		UserRpc:         user.NewUser(zrpc.MustNewClient(c.UserRpc, rpcOpt)),
		NotificationRpc: notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc, rpcOpt)),
		OpenRpc:         open.NewOpen(zrpc.MustNewClient(c.OpenRpc, rpcOpt)),
		RocketMQ:        corerocketmq.InitRocketMQ(c.RocketMQ.Addr),
	}
}
