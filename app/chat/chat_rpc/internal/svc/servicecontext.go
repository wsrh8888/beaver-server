package svc

import (
	"beaver/app/chat/chat_rpc/internal/config"
	"beaver/app/user/user_rpc/user"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"
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
	RocketMQ   *corerocketmq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)

	userRpc := user.NewUser(zrpc.MustNewClient(c.UserRpc))

	// 初始化 RocketMQ 客户端

	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)

	return &ServiceContext{
		Config:     c,
		Redis:      client,
		DB:         mysqlDb,
		VersionGen: versionGen,
		UserRpc:    userRpc,
		RocketMQ:   mqClient,
	}
}
