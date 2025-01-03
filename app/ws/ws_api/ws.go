package main

import (
	"flag"
	"fmt"

	"beaver/app/ws/ws_api/internal/config"
	"beaver/app/ws/ws_api/internal/handler"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/common/etcd"
	middleware "beaver/common/middle"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ws.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	server.Use(middleware.LogMiddleware)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
