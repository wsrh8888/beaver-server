package config

import "github.com/zeromicro/go-zero/zrpc"

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
	ApiBaseUrl string // Webhook URL 基址，如 https://api.example.com
	UserRpc    zrpc.RpcClientConf
}
