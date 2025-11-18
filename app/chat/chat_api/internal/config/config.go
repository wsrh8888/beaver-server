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
	UserRpc   zrpc.RpcClientConf
	ChatRpc   zrpc.RpcClientConf
	FriendRpc zrpc.RpcClientConf
}
