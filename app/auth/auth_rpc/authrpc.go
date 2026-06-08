package main

import (
	"flag"
	"fmt"

	"beaver/app/auth/auth_rpc/internal/config"
	"beaver/app/auth/auth_rpc/internal/server"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/utils/logger"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/authrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logger.Init("auth_rpc")
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		auth_rpc.RegisterAuthServer(grpcServer, server.NewAuthServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
