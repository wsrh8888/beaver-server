package svc

import (
	"beaver/app/user/user_rpc/internal/config"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *redis.Client
	VersionGen *versionPkg.VersionGenerator
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	return &ServiceContext{
		Config:     c,
		Redis:      client,
		DB:         mysqlDb,
		VersionGen: versionGen,
	}
}
