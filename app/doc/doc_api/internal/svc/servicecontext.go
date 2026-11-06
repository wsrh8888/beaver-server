package svc

import (
	"beaver/app/doc/doc_api/internal/config"
	doccli "beaver/app/doc/doc_rpc/doc"
	"beaver/common/zrpc_interceptor"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	Doc    doccli.Doc
}

func NewServiceContext(c config.Config) *ServiceContext {
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	return &ServiceContext{
		Config: c,
		Doc:    doccli.NewDoc(zrpc.MustNewClient(c.DocRpc, rpcOpt)),
	}
}
