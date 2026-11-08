package ymsgw

import (
	"context"
	"net/http"
	"time"

	"beaver/app/doc/doc_api/internal/config"
	doccli "beaver/app/doc/doc_rpc/doc"

	"github.com/go-redis/redis"
	goredis "github.com/redis/go-redis/v9"
	ygoredis "github.com/reearth/ygo/cluster/redis"
	ygws "github.com/reearth/ygo/provider/websocket"
)

// Gateway 封装 ygo websocket.Server 及其持久化、鉴权、跨实例广播依赖，供 doc_api 复用。
type Gateway struct {
	srv   *ygws.Server
	relay *ygoredis.Relay
}

// New 构建协同网关：注入 doc_rpc 持久化适配器、升级前鉴权钩子、Redis 跨实例房间广播。
func New(c config.Config, doc doccli.Doc, rdbV6 *redis.Client, rdbV9 *goredis.Client) (*Gateway, error) {
	adapter := &persistence{doc: doc}
	srv := ygws.NewServerWithPersistence(adapter)

	// 反向绑定：Compact 需读取房间内实时 doc 解码正文。
	adapter.srv = srv

	a := &wsAuth{secret: c.Auth.AccessSecret, rdb: rdbV6, doc: doc}
	srv.Authorize = a.authorize

	if c.Yjs.CompactEvery > 0 {
		srv.CompactEvery = c.Yjs.CompactEvery
	}
	if c.Yjs.AutoVersionEvery > 0 {
		srv.AutoVersionEvery = time.Duration(c.Yjs.AutoVersionEvery) * time.Second
	}
	if c.Yjs.RoomIdleTimeout > 0 {
		srv.RoomIdleTimeout = time.Duration(c.Yjs.RoomIdleTimeout) * time.Second
	}
	if c.Yjs.MaxPeersPerRoom > 0 {
		srv.MaxPeersPerRoom = c.Yjs.MaxPeersPerRoom
	}
	srv.HocuspocusFraming = c.Yjs.HocuspocusFraming

	// Redis 跨实例房间广播：AttachRelay 内部完成 relay.Start。
	relay, err := ygoredis.New(rdbV9, ygoredis.Config{ChannelPrefix: c.Yjs.RelayChannelPrefix})
	if err != nil {
		return nil, err
	}
	if err := srv.AttachRelay(relay); err != nil {
		return nil, err
	}

	return &Gateway{srv: srv, relay: relay}, nil
}

// ServeHTTP 处理 WebSocket 升级及后续 sync/awareness 帧。调用方需先把 r.URL.Path 置为
// "/"+docId，使 ygo 以 docId 作为房间名。
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.srv.ServeHTTP(w, r)
}

// Shutdown 优雅关闭：停止接收新连接、卸载房间并触发末次压实落库。
func (g *Gateway) Shutdown(ctx context.Context) error {
	return g.srv.Shutdown(ctx)
}
