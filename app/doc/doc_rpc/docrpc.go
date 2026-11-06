package main

import (
	"flag"
	"fmt"

	"beaver/app/doc/doc_rpc/internal/config"
	"beaver/app/doc/doc_rpc/internal/server"
	"beaver/app/doc/doc_rpc/internal/svc"
	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"
	"beaver/utils/beaverlog"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/docrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	beaverlog.InitFromConf(c.RpcServerConf.ServiceConf)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		doc_rpc.RegisterDocServer(grpcServer, server.NewDocServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
