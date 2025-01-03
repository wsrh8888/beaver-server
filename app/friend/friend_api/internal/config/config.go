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
	Etcd  string
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	UserRpc   zrpc.RpcClientConf
	FriendRpc zrpc.RpcClientConf
	ChatRpc   zrpc.RpcClientConf
}
