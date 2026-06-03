package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	Auth struct {
		AccessSecret string
	}
	Etcd     string
	GroupRpc zrpc.RpcClientConf
	ChatRpc  zrpc.RpcClientConf
	RocketMQ struct {
		Addr string
	}
	InstanceID string `json:",optional"`
	WebSocket struct {
		PongWait             int
		WriteWait            int
		PingPeriod           int
		MaxMessageSize       int
		AppHeartbeatInterval int
	}
}
