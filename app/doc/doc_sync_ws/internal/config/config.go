package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	Etcd       string
	DocRpc     zrpc.RpcClientConf
	InstanceID string `json:",optional"`
	// Yjs 协同网关参数（映射 ygo websocket.Server 选项）。
	Yjs struct {
		CompactEvery       int    `json:",optional"`               // 累计多少条增量后压实快照，<=0 用库默认
		AutoVersionEvery   int    `json:",optional"`               // 秒；>0 周期性生成用户可见版本，0 关闭
		RoomIdleTimeout    int    `json:",optional"`               // 秒；房间空闲多久后卸载
		MaxPeersPerRoom    int    `json:",optional"`               // 单房间最大连接数，0 不限
		RelayChannelPrefix string `json:",optional"`               // Redis 跨实例广播 channel 前缀
		HocuspocusFraming  bool   `json:",optional,default=false"` // 是否启用 Hocuspocus docName 前缀帧
	}
}
