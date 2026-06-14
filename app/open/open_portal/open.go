package main

import (
	"beaver/app/open/open_portal/internal/config"
	"beaver/app/open/open_portal/internal/handler"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/common/etcd"
	"beaver/utils/logger"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/open.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logger.Init("open_portal")

	// 配置 CORS - 允许所有来源和头部
	server := rest.MustNewServer(c.RestConf,
		rest.WithCors("*"),
		rest.WithCorsHeaders("*"),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
