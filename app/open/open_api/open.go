package main

import (
	"beaver/app/open/open_api/internal/config"
	"beaver/app/open/open_api/internal/handler"
	"beaver/app/open/open_api/internal/svc"
	"beaver/common/etcd"
	commonMiddleware "beaver/common/middleware/http"
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

	logger.Init("open_api")

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册全局中间件（日志）
	server.Use(commonMiddleware.RequestLogMiddleware)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
