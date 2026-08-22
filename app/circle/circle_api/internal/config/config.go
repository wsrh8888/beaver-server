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
	Etcd            string
	Domain          string // 对外域名，如 http://192.168.3.4:20800
	RocketMQ        struct {
		Addr string
	}
	UserRpc         zrpc.RpcClientConf
	ChatRpc         zrpc.RpcClientConf
	NotificationRpc zrpc.RpcClientConf
}
