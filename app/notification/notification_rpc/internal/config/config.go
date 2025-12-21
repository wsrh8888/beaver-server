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
	UserRpc zrpc.RpcClientConf
}
