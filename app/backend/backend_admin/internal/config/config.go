package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd  string
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	WhiteList []string //白名单
	UserRpc   zrpc.RpcClientConf
	Auth      struct {
		AccessSecret string
		AccessExpire int
	}
	FileMaxSize map[string]float64
	BlackList   []string
	UploadDir   string
	FileRpc     zrpc.RpcClientConf
	Qiniu       struct {
		AK         string
		SK         string
		Bucket     string
		Domain     string
		ExpireTime int64
	}
	DictionaryRpc zrpc.RpcClientConf
}
