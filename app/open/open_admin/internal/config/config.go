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
	Etcd string
	JWT  struct {
		SecretKey   string
		ExpireHours int
	}
	UserRpc zrpc.RpcClientConf
}
