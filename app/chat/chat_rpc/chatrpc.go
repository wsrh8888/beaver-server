package main

import (
	"flag"
	"fmt"

	"beaver/app/chat/chat_rpc/internal/config"
	"beaver/app/chat/chat_rpc/internal/server"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	grpcMiddleware "beaver/common/middle/grpc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/chatrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		chat_rpc.RegisterChatServer(grpcServer, server.NewChatServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 添加拦截器
	s.AddUnaryInterceptors(grpcMiddleware.RequestLogInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
