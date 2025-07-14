package svc

import (
	"beaver/app/dictionary/dictionary_api/internal/config"
	"beaver/app/dictionary/dictionary_rpc/types/dictionary_rpc"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	DictionaryRpc dictionary_rpc.DictionaryClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		DictionaryRpc: dictionary_rpc.NewDictionaryClient(zrpc.MustNewClient(c.DictionaryRpc).Conn()),
	}
}
