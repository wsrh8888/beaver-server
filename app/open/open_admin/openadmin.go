package main

import (
	"beaver/app/open/open_admin/internal/config"
	"beaver/app/open/open_admin/internal/handler"
	"beaver/app/open/open_admin/internal/middleware"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/common/etcd"
	commonMiddleware "beaver/common/middleware"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/openadmin.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册全局中间件
	server.Use(commonMiddleware.RequestLogMiddleware)

	// JWT 鉴权中间件（排除登录接口）
	authMiddleware := middleware.AdminAuthMiddleware(c.JWT.SecretKey)
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 登录接口跳过鉴权
			if strings.HasPrefix(r.URL.Path, "/admin/open/auth/login") {
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
