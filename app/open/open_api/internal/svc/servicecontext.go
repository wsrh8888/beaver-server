package svc

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/config"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core/coregorm"
	"time"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	UserRpc user_rpc.UserClient
	ChatRpc chat_rpc.ChatClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:  c,
		DB:      db,
		UserRpc: user_rpc.NewUserClient(zrpc.MustNewClient(c.UserRpc, zrpc.WithTimeout(time.Duration(c.UserRpc.Timeout)*time.Millisecond)).Conn()),
		ChatRpc: chat_rpc.NewChatClient(zrpc.MustNewClient(c.ChatRpc, zrpc.WithTimeout(time.Duration(c.ChatRpc.Timeout)*time.Millisecond)).Conn()),
	}
}
