// Package circle A 主体（开源核心）→ 收费模块圈子服务（beaver-server-social/app/circle）的 S2S 接口代理。
//
// 与 doc / moment 同构，但域名、凭据、可调接口都只属于这一端：一个是圈子在调，
// 一个是动态在调，混在一起会变成「谁的票都能调谁的接口」。
//
// 装配在 servicecontext 构造时：NewCircle 灌入凭据与域名，域名写入 commonconfig.CircleUrl。
// 业务侧经 app/circle/circle_rpc/circle 的接口调用，最终落到本包的 Post/Get/Delete。
package circle

import (
	"beaver/commonconfig"
	"sync"
	"time"
)

// client 本包唯一的出站客户端，Init 时装配。
var client *Client

// Init 装配出站凭据（收费模块圈子分配给 A 主体的 client_credentials）与域名。
// 只在该服务装配依赖时调用一次；域名写入 commonconfig.CircleUrl。
func Init(clientID, clientSecret, baseURL string) {
	commonconfig.CircleUrl = baseURL
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
