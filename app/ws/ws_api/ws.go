package main

import (
	"flag"
	"fmt"

	"beaver/app/ws/ws_api/internal/config"
	"beaver/app/ws/ws_api/internal/handler"
	"beaver/app/ws/ws_api/internal/logic"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/common/etcd"
	"beaver/common/middleware"
	"beaver/utils/logger"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ws.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logger.Init("ws_api")

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 启动 RocketMQ Consumer
	if ctx.RocketMQ != nil {
		mqConsumer := logic.NewMqConsumerLogic(ctx)
		go func() {
			if err := mqConsumer.StartConsumer(); err != nil {
				fmt.Printf("RocketMQ Consumer 启动失败: %v\n", err)
			}
		}()
	}

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	server.Use(middleware.LogMiddleware)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
