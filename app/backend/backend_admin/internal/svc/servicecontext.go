package svc

import (
	"beaver/app/backend/backend_admin/internal/config"
	"beaver/app/dictionary/dictionary_rpc/dictionary"
	"beaver/app/dictionary/dictionary_rpc/types/dictionary_rpc"
	"beaver/app/file/file_rpc/file"
	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	Redis         *redis.Client
	UserRpc       user_rpc.UserClient
	FileRpc       file_rpc.FileClient
	DictionaryRpc dictionary_rpc.DictionaryClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)

	return &ServiceContext{
		Config:        c,
		DB:            mysqlDb,
		Redis:         client,
		FileRpc:       file.NewFile(zrpc.MustNewClient(c.FileRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		DictionaryRpc: dictionary.NewDictionary(zrpc.MustNewClient(c.DictionaryRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		UserRpc:       user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
