package svc

import (
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/config"
	"beaver/app/open/open_api/internal/middleware"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/core/coregorm"
	"beaver/core/corewebhook"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	UserRpc        user_rpc.UserClient
	ChatRpc        chat_rpc.ChatClient
	WebhookSender  *corewebhook.WebhookSender
	AuthMiddleware func(http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	return &ServiceContext{
		Config:         c,
		DB:             db,
		UserRpc:        user_rpc.NewUserClient(zrpc.MustNewClient(c.UserRpc, zrpc.WithTimeout(time.Duration(c.UserRpc.Timeout)*time.Millisecond)).Conn()),
		ChatRpc:        chat_rpc.NewChatClient(zrpc.MustNewClient(c.ChatRpc, zrpc.WithTimeout(time.Duration(c.ChatRpc.Timeout)*time.Millisecond)).Conn()),
		WebhookSender:  corewebhook.NewWebhookSender(db, corewebhook.Config{Timeout: 10, RetryCount: 3}),
		AuthMiddleware: middleware.AuthMiddleware(db),
	}
}
