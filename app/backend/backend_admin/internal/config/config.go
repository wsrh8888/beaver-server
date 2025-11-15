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
	BlackList []string //黑名单
	File      struct {
		WhiteList []string
		BlackList []string
		MaxSize   map[string]float64 // 文件大小限制
	}
	UserRpc zrpc.RpcClientConf
	Auth    struct {
		AccessSecret string
		AccessExpire int
	}
	Local struct {
		UploadDir string // 本地文件上传目录
	}
	Qiniu struct {
		AK         string
		SK         string
		Bucket     string
		Domain     string
		ExpireTime int64 // 签名URL有效期, 单位：秒
	}
	FileRpc       zrpc.RpcClientConf
	DictionaryRpc zrpc.RpcClientConf
}
