package svc

import (
	"beaver/app/datasync/datasync_api/internal/config"
	"beaver/app/datasync/datasync_rpc/types/types/datasync_rpc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *redis.Client
	DatasyncRpc datasync_rpc.DatasyncClient
	FriendRpc   friend_rpc.FriendClient
	GroupRpc    group_rpc.GroupClient
	UserRpc     user_rpc.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:      c,
		DB:          mysqlDb,
		Redis:       client,
		DatasyncRpc: datasync_rpc.NewDatasyncClient(zrpc.MustNewClient(c.DatasyncRpc).Conn()),
		FriendRpc:   friend_rpc.NewFriendClient(zrpc.MustNewClient(c.FriendRpc).Conn()),
		GroupRpc:    group_rpc.NewGroupClient(zrpc.MustNewClient(c.GroupRpc).Conn()),
		UserRpc:     user_rpc.NewUserClient(zrpc.MustNewClient(c.UserRpc).Conn()),
	}
}
