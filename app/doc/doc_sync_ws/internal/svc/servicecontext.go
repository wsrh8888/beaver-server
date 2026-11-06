package svc

import (
	doccli "beaver/app/doc/doc_rpc/doc"
	"beaver/app/doc/doc_sync_ws/internal/config"
	"beaver/app/doc/doc_sync_ws/internal/ymsgw"
	"beaver/common/zrpc_interceptor"
	"beaver/core/coreredis"

	"github.com/go-redis/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	Doc     doccli.Doc
	Redis   *redis.Client
	Gateway *ymsgw.Gateway
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 登录态校验用 go-redis v6（与 ws_api 一致）。
	rdbV6 := coreredis.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.Db)

	// doc_rpc 权威存储客户端。
	rpcOpt := zrpc.WithUnaryClientInterceptor(zrpc_interceptor.ClientInfoInterceptor)
	doc := doccli.NewDoc(zrpc.MustNewClient(c.DocRpc, rpcOpt))

	// ygo 跨实例房间广播用 go-redis v9。
	rdbV9 := goredis.NewClient(&goredis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.Db,
	})

	gw, err := ymsgw.New(c, doc, rdbV6, rdbV9)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:  c,
		Doc:     doc,
		Redis:   rdbV6,
		Gateway: gw,
	}
}
