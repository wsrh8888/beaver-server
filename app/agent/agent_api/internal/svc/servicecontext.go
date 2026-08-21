package svc

import (
	"beaver/app/agent/agent_api/internal/config"
	"beaver/app/agent/agent_rpc/agent"
	"beaver/app/agent/pyagent"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	AgentRpc    agent.Agent
	BeaverAgent pyagent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	return &ServiceContext{
		Config:   c,
		DB:       coregorm.InitGorm(c.Mysql.DataSource),
		AgentRpc: agent.NewAgent(zrpc.MustNewClient(c.AgentRpc, rpcOpt)),
		BeaverAgent: pyagent.NewClient(zrpc.MustNewClient(zrpc.RpcClientConf{
			Endpoints: c.BeaverAgent.Endpoints,
			NonBlock:  c.BeaverAgent.NonBlock,
		}, rpcOpt)),
	}
}
