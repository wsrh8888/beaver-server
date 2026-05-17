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
	Etcd     string
	UserRpc  zrpc.RpcClientConf
	ChatRpc  zrpc.RpcClientConf
	GroupRpc zrpc.RpcClientConf
	OAuthRpc zrpc.RpcClientConf // OAuth RPC 服务
	OAuth    OAuthConf          // OAuth 授权页面配置
}

// OAuthConf OAuth 授权页面配置
type OAuthConf struct {
	BaseUrl string // OAuth 服务基础 URL（授权页面地址）
}
