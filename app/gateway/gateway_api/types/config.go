package types

import "github.com/zeromicro/go-zero/core/logx"

type Config struct {
	Name       string
	Addr       string `json:",optional"`
	Etcd       string
	Log        logx.LogConf
	Prometheus PrometheusConfig
	Limit          LimitConfig
	Auth           AuthConfig
	PublicList     []string `json:",optional"` // Gateway 不鉴权（含 *_public）
	CustomAuthList []string `json:",optional"` // 透传，由下游服务 middleware 鉴权
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int
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
