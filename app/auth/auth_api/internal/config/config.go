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
	WhiteList []string //白名单
	UserRpc   zrpc.RpcClientConf
	Auth      struct {
		AccessSecret string
		AccessExpire int
	}
	Email struct {
		QQ struct {
			Host     string
			Port     int
			Username string
			Password string
		}
	}
}
