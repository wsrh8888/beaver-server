package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd string
	// DocRpc：文档权威域 gRPC（space/node/doc/perm/revision 落库均经 doc_rpc）。
	DocRpc zrpc.RpcClientConf
}
