package main

import (
	"beaver/app/open/open_portal/internal/config"
	"beaver/app/open/open_portal/internal/handler"
	"beaver/app/open/open_portal/internal/middleware"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/common/etcd"
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

	// 配置 CORS - 允许所有来源和头部
	server := rest.MustNewServer(c.RestConf,
		rest.WithCors("*"),
		rest.WithCorsHeaders("*"),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册全局中间件（日志）

	// 开发者认证中间件（所有接口都需要登录，但部分接口不需要开发者资质）
	devAuthMiddleware := middleware.DeveloperAuthMiddleware(c.Auth.AccessSecret, ctx.DB)
	server.Use(devAuthMiddleware)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
