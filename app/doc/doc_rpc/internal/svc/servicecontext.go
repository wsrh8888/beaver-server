package svc

import (
	"beaver/app/doc/doc_rpc/internal/config"
	knowledgecli "beaver/app/knowledge/knowledge_rpc/knowledge"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coregorm"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	KnowledgeRpc knowledgecli.Knowledge
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := coregorm.InitGorm(c.Mysql.DataSource)
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	return &ServiceContext{
		Config:       c,
		DB:           db,
		KnowledgeRpc: knowledgecli.NewKnowledge(zrpc.MustNewClient(c.KnowledgeRpc, rpcOpt)),
	}
}
