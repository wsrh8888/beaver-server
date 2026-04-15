package main

import (
	"beaver/app/open/open_api/internal/config"
	"beaver/app/open/open_api/internal/handler"
	"beaver/app/open/open_api/internal/middleware"
	"beaver/app/open/open_api/internal/svc"
	"beaver/common/etcd"
	commonMiddleware "beaver/common/middleware"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/open.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册全局中间件（日志）
	server.Use(commonMiddleware.RequestLogMiddleware)

	// OAuth2 鉴权中间件（排除认证接口）
	authMiddleware := middleware.AuthMiddleware(ctx.DB)
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 认证接口跳过鉴权
			if strings.HasPrefix(r.URL.Path, "/api/open/v1/auth/") {
				next(w, r)
				return
			}
			// 其他接口需要鉴权
			authMiddleware(next)(w, r)
		}
	})

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
