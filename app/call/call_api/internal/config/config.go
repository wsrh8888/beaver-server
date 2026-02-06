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
		DB       int
	}
	LiveKit struct {
		Host      string
		ApiKey    string
		ApiSecret string
	}
	UserRpc zrpc.RpcClientConf
	CallRpc zrpc.RpcClientConf
}
