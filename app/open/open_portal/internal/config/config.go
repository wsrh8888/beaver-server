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
	OAuth   struct {
		BaseUrl string // OAuth 服务基础 URL（授权页面地址）
	}
}

// OAuthConf OAuth 授权页面配置
