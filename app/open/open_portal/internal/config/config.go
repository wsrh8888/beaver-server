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
	UserRpc  zrpc.RpcClientConf
	AuthRpc  zrpc.RpcClientConf
	GroupRpc zrpc.RpcClientConf
	OpenRpc  zrpc.RpcClientConf
	OAuth    struct {
		BaseUrl string // OAuth 授权页基础 URL（beaver-oauth H5）
	}
	Domain      string // 对外 API 网关地址（webhook、code 换 token 等）
	PortalOAuth struct {
		AppId string // 开放平台门户自身作为 OAuth 应用的 AppId
	}
}
