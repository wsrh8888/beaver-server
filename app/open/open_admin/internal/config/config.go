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
	Auth struct {
		AccessSecret string
		AccessExpire int
	}
	UserRpc zrpc.RpcClientConf
}
