package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd string
	Mysql struct {
		DataSource string
	}
	// AgentRpc：Agent CRUD / 历史
	AgentRpc zrpc.RpcClientConf
	// BeaverAgent：Python 流式大脑（直连）
	BeaverAgent struct {
		Endpoints []string
		NonBlock  bool `json:",default=true"`
	}
}
