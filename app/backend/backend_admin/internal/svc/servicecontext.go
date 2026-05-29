package svc

import (
	"beaver/app/backend/backend_admin/internal/config"
	"beaver/app/file/file_rpc/file"
	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/app/user/user_rpc/user"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *redis.Client
	VersionGen *versionPkg.VersionGenerator
	UserRpc    user_rpc.UserClient
	FileRpc    file_rpc.FileClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	return &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		Redis:      client,
		VersionGen: versionGen,
		FileRpc:    file.NewFile(zrpc.MustNewClient(c.FileRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
		UserRpc:    user.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor))),
	}
}
