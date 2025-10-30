package main

import (
	"flag"
	"fmt"

	"beaver/app/datasync/datasync_rpc/internal/config"
	"beaver/app/datasync/datasync_rpc/internal/server"
	"beaver/app/datasync/datasync_rpc/internal/svc"
	"beaver/app/datasync/datasync_rpc/types/types/datasync_rpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/datasyncrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		datasync_rpc.RegisterDatasyncServer(grpcServer, server.NewDatasyncServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
