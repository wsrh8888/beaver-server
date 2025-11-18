package main

import (
	"beaver/app/update/update_api/internal/config"
	"beaver/app/update/update_api/internal/handler"
	"beaver/app/update/update_api/internal/svc"
	"beaver/common/etcd"
	"beaver/common/middleware"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/update.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.Use(middleware.LogMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
