package svc

import (
	"beaver/app/notification/notification_rpc/internal/config"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core"
	versionPkg "beaver/core/version"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user_rpc.UserClient
	Redis   *redis.Client

	DB         *gorm.DB
	VersionGen *versionPkg.VersionGenerator
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)

	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	return &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		VersionGen: versionGen,
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
