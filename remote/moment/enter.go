// Package moment A 主体（开源核心）→ 圈子动态服务（beaver-server-social/app/social）的 S2S 接口代理。
//
// 跨仓不允许 etcd gRPC 直连（团队铁律），统一走对方标准路径 /beaver-social/api/moment/machine/v1/*（client_credentials 换票 + Bearer）。
// 结构对齐 remote/{open,ai,circle,doc}：enter.go 灌域名与凭据、remote.go 持有换票与发包 plumbing、
// 业务方法在 moment.go。Req/Res 放到 types/remote/moment。
package moment

import (
	"sync"
	"time"
)

// client 本包唯一的出站客户端，Init 时装配。
var client *Client

// Init 装配出站凭据（收费模块动态服务分配给 A 主体的 client_credentials）与域名。
// 只在该服务装配依赖时调用一次；域名写入 commonconfig.SocialUrl。
func Init(clientID, clientSecret string) {
	client = &Client{clientID: clientID, clientSecret: clientSecret}
}

// Client 持凭据与换来的票，换票缓存挂在它上面。
type Client struct {
	clientID     string
	clientSecret string

	mu        sync.Mutex
	bearer    string
	expiresAt time.Time
}
