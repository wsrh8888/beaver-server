package svc

import (
	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/datasync/datasync_rpc/internal/config"
	"beaver/app/friend/friend_rpc/friend"
	"beaver/app/group/group_rpc/group"
	"beaver/app/user/user_rpc/user"
	"beaver/core"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	DB        *gorm.DB
	Redis     *redis.Client
	ChatRpc   chat.Chat
	FriendRpc friend.Friend
	GroupRpc  group.Group
	UserRpc   user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)
	chatRpc := chat.NewChat(zrpc.MustNewClient(c.ChatRpc))
	friendRpc := friend.NewFriend(zrpc.MustNewClient(c.FriendRpc))
	groupRpc := group.NewGroup(zrpc.MustNewClient(c.GroupRpc))
	userRpc := user.NewUser(zrpc.MustNewClient(c.UserRpc))

	return &ServiceContext{
		Config:    c,
		Redis:     client,
		DB:        mysqlDb,
		ChatRpc:   chatRpc,
		FriendRpc: friendRpc,
		GroupRpc:  groupRpc,
		UserRpc:   userRpc,
	}
}
