package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd string

	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	FriendRpc zrpc.RpcClientConf
	GroupRpc  zrpc.RpcClientConf
	UserRpc   zrpc.RpcClientConf
	ChatRpc   zrpc.RpcClientConf
	MomentRpc zrpc.RpcClientConf
	EmojiRpc  zrpc.RpcClientConf
}
