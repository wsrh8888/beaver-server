package svc

import (
	"beaver/app/emoji/emoji_rpc/internal/config"
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

	return &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		Redis:      client,
		VersionGen: versionPkg.NewVersionGenerator(client, mysqlDb),
	}
}
