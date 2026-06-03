package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	RedisConf struct {
		Addr     string
		Password string
		Db       int
	}
	RocketMQ struct {
		Addr string
	}
	DatasyncRpc zrpc.RpcClientConf
	UserRpc     zrpc.RpcClientConf
	NotificationRpc zrpc.RpcClientConf
	Push        struct {
		Enabled bool
		FCM     struct {
			Enabled         bool
			ProjectID       string
			CredentialsFile string
		}
		APNs struct {
			Enabled    bool
			KeyFile    string
			KeyID      string
			TeamID     string
			BundleID   string
			Production bool
		}
	}
}
