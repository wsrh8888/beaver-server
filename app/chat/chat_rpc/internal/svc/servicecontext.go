package svc

import (
	"beaver/app/chat/chat_rpc/internal/config"
	"beaver/app/notification/notification_rpc/notification"
	"beaver/app/user/user_rpc/user"
	"beaver/core/corepush"
	"beaver/core/coregorm"
	"beaver/core/coreredis"
	"beaver/core/corerocketmq"
	"beaver/core/corewebhook"
	versionPkg "beaver/core/version"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	Redis         *redis.Client
	VersionGen    *versionPkg.VersionGenerator
	UserRpc       user.User
	NotificationRpc notification.Notification
	RocketMQ      *corerocketmq.Client
	WebhookSender *corewebhook.WebhookSender
	PushSender    *corepush.PushSender
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := coregorm.InitGorm(c.Mysql.DataSource)
	client := coreredis.InitRedis(c.RedisConf.Addr, c.RedisConf.Password, c.RedisConf.Db)
	versionGen := versionPkg.NewVersionGenerator(client, mysqlDb)
	userRpc := user.NewUser(zrpc.MustNewClient(c.UserRpc))
	notificationRpc := notification.NewNotification(zrpc.MustNewClient(c.NotificationRpc))
	mqClient := corerocketmq.InitRocketMQ(c.RocketMQ.Addr)

	return &ServiceContext{
		Config:        c,
		Redis:         client,
		DB:            mysqlDb,
		VersionGen:    versionGen,
		UserRpc:       userRpc,
		NotificationRpc: notificationRpc,
		RocketMQ:      mqClient,
		WebhookSender: corewebhook.NewWebhookSender(mysqlDb, corewebhook.Config{Timeout: 10, RetryCount: 3}),
		PushSender: corepush.NewPushSender(corepush.Config{
			Enabled: c.Push.Enabled,
			FCM: corepush.FCMConfig{
				Enabled:         c.Push.FCM.Enabled,
				ProjectID:       c.Push.FCM.ProjectID,
				CredentialsFile: c.Push.FCM.CredentialsFile,
			},
			APNs: corepush.APNsConfig{
				Enabled:    c.Push.APNs.Enabled,
				KeyFile:    c.Push.APNs.KeyFile,
				KeyID:      c.Push.APNs.KeyID,
				TeamID:     c.Push.APNs.TeamID,
				BundleID:   c.Push.APNs.BundleID,
				Production: c.Push.APNs.Production,
			},
		}),
	}
}
