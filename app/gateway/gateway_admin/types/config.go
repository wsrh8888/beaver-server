package types

import "github.com/zeromicro/go-zero/core/logx"

type Config struct {
	Name string
	Addr string `json:",optional"`
	Etcd string
	Log  logx.LogConf
}
