package svc

import (
	"beaver/app/circle/circle_rpc/internal/config"
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
	redisClient := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)
	return &ServiceContext{
		Config:     c,
		DB:         mysqlDb,
		Redis:      redisClient,
		VersionGen: versionPkg.NewVersionGenerator(redisClient, mysqlDb),
	}
}
