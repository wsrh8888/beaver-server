package types

import "github.com/zeromicro/go-zero/core/logx"

type Config struct {
	Addr string
	Etcd string
	Log  logx.LogConf
}
