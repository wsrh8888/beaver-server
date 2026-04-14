package svc

import (
	"beaver/app/call/call_rpc/internal/config"
	"beaver/core/coregorm"
	"beaver/core/coreredis"

	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/user/user_rpc/user"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redis.Client
	ChatRpc chat.Chat
	UserRpc user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	rdb := coreredis.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)

	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   rdb,
		ChatRpc: chat.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		UserRpc: user.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
