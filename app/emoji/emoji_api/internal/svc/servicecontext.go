package svc

import (
	"beaver/app/emoji/emoji_api/internal/config"
	"beaver/core"
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
	mysqlDb := core.InitGorm(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	return &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		Redis:      client,
		VersionGen: versionGen,
	}
}
