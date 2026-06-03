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
	AuthRpc zrpc.RpcClientConf
	OAuth   struct {
		BaseUrl string // OAuth 服务基础 URL（授权页面地址）
	}
	PortalOAuth struct {
		AppId      string // 开放平台门户自身作为 OAuth 应用的 AppId
		GatewayUrl string // API 网关地址，用于 code 换 token
	}
	ApiBaseUrl string // open_api 对外根地址，用于拼接 Incoming Webhook URL
}

// OAuthConf OAuth 授权页面配置
