package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Domain string // 对外访问域名（用于生成本地文件完整 URL）
	Etcd   string
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
	AuthRpc zrpc.RpcClientConf
	Auth    struct {
		AccessSecret string
		AccessExpire int
	}
	Local struct {
		UploadDir   string // 本地文件上传目录
		ProjectName string // 项目名称，用于文件路径前缀
	}
	Qiniu struct {
		ProjectName string // 项目名称，用于文件路径前缀
		AK          string
		SK          string
		Bucket      string
		Domain      string
		ExpireTime  int64 // 签名URL有效期, 单位：秒
	}
	FileRpc      zrpc.RpcClientConf
	PlatformRpc  zrpc.RpcClientConf
	OpenRpc      zrpc.RpcClientConf
	FriendRpc    zrpc.RpcClientConf
	GroupRpc     zrpc.RpcClientConf
	ChatRpc      zrpc.RpcClientConf
	MomentRpc    zrpc.RpcClientConf
	EmojiRpc     zrpc.RpcClientConf
}
