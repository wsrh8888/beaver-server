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
	FileMaxSize map[string]float64
	WhiteList   []string
	BlackList   []string
	UserRpc     zrpc.RpcClientConf
	Local       struct {
		UploadDir string
	}
	Qiniu struct {
		AK         string
		SK         string
		Bucket     string
		Domain     string
		ExpireTime int64
	}
}
