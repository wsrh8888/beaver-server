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
	Domain      string // 项目对外访问域名（用于生成本地文件完整URL）
	FileMaxSize map[string]float64
	WhiteList   []string
	BlackList   []string
	UserRpc     zrpc.RpcClientConf
	Local       struct {
		UploadDir   string // 本地文件上传目录
		ProjectName string // 项目名称，用于文件路径前缀（为空则使用根目录）
	}
	Qiniu struct {
		ProjectName string // 项目名称，用于文件路径前缀（为空则使用根目录）
		AK          string
		SK          string
		Bucket      string
		Domain      string // 七牛云文件访问域名
		ExpireTime  int64
	}
}
