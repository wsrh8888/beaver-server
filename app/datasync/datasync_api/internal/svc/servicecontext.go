package svc

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_api/internal/config"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core"
	"time"

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
	EmojiRpc  emoji_rpc.EmojiClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:    c,
		DB:        mysqlDb,
		Redis:     client,
		FriendRpc: friend_rpc.NewFriendClient(zrpc.MustNewClient(c.FriendRpc, zrpc.WithTimeout(time.Duration(c.FriendRpc.Timeout)*time.Millisecond)).Conn()),
		GroupRpc:  group_rpc.NewGroupClient(zrpc.MustNewClient(c.GroupRpc, zrpc.WithTimeout(time.Duration(c.GroupRpc.Timeout)*time.Millisecond)).Conn()),
		UserRpc:   user_rpc.NewUserClient(zrpc.MustNewClient(c.UserRpc, zrpc.WithTimeout(time.Duration(c.UserRpc.Timeout)*time.Millisecond)).Conn()),
		ChatRpc:   chat_rpc.NewChatClient(zrpc.MustNewClient(c.ChatRpc, zrpc.WithTimeout(time.Duration(c.ChatRpc.Timeout)*time.Millisecond)).Conn()),
		EmojiRpc:  emoji_rpc.NewEmojiClient(zrpc.MustNewClient(c.EmojiRpc, zrpc.WithTimeout(time.Duration(c.EmojiRpc.Timeout)*time.Millisecond)).Conn()),
	}
}
