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
	Etcd      string
	GroupRpc  zrpc.RpcClientConf
	WebSocket struct {
		PongWait             int
		WriteWait            int
		PingPeriod           int
		MaxMessageSize       int
		AppHeartbeatInterval int
	}
}
