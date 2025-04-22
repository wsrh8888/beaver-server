package types

import "github.com/zeromicro/go-zero/core/logx"

type Config struct {
	Name       string
	Addr       string `json:",optional"`
	Etcd       string
	Log        logx.LogConf
	Prometheus PrometheusConfig
	Limit      LimitConfig
}

type PrometheusConfig struct {
	Enable bool   `json:",default=false"`
	Path   string `json:",default=/metrics"`
}

type LimitConfig struct {
	Enable bool    `json:",default=false"`
	Rate   float64 `json:",default=100"`
	Burst  int     `json:",default=200"`
}
