// Package doc A 主体（开源核心）→ 在线文档服务（beaver-doc 独立仓）的 S2S 接口代理。
//
// 跨仓不允许 etcd gRPC 直连（团队铁律），统一走对方网关 /internal/Doc.*（client_credentials 换票 + Bearer）。
// 结构对齐 remote/{open,ai,circle,moment}：enter.go 灌域名与凭据、remote.go 持有换票与发包 plumbing、
// 业务方法在 content.go。Req/Res 放到 types/remote/doc。
package doc

import (
	"beaver/commonconfig"
	"sync"
	"time"
)

// client 本包唯一的出站客户端，Init 时装配。
var client *Client

// Init 装配出站凭据（beaver-doc 分配给调用方的 client_credentials）与域名。
// 只在该服务装配依赖时调用一次；域名写入 commonconfig.DocUrl。
func Init(clientID, clientSecret, baseURL string) {
	commonconfig.DocUrl = baseURL
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
