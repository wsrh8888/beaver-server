package svc

import (
	"beaver/app/chat/chat_rpc/chat"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/document/document_api/internal/config"
	"beaver/app/file/file_rpc/file"
	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redis.Client
	ChatRpc chat_rpc.ChatClient
	FileRpc file_rpc.FileClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   client,
		ChatRpc: chat.NewChat(zrpc.MustNewClient(c.ChatRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		FileRpc: file.NewFile(zrpc.MustNewClient(c.FileRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
