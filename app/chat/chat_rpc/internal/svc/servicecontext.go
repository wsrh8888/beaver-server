package svc

import (
	"beaver/app/chat/chat_rpc/internal/config"
	"beaver/app/user/user_rpc/user"
	"beaver/core"
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
	UserRpc    user.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	userRpc := user.NewUser(zrpc.MustNewClient(c.UserRpc))

	return &ServiceContext{
		Config:     c,
		Redis:      client,
		DB:         mysqlDb,
		VersionGen: versionGen,
		UserRpc:    userRpc,
	}
}
