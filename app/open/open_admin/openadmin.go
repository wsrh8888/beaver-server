package main

import (
	"beaver/app/open/open_admin/internal/config"
	"beaver/app/open/open_admin/internal/handler"
	"beaver/app/open/open_admin/internal/middleware"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/common/etcd"
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

	// 注册全局中间件（日志）

	// 开发者认证中间件（排除登录和申请接口）
	devAuthMiddleware := middleware.DeveloperAuthMiddleware(c.Auth.AccessSecret, ctx.DB)
	server.Use(func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 白名单接口跳过鉴权
			if isWhiteListPath(r.URL.Path) {
				next(w, r)
				return
			}
			// 其他接口需要开发者认证
			devAuthMiddleware(next)(w, r)
		}
	})

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// isWhiteListPath 判断是否是白名单路径
func isWhiteListPath(path string) bool {
	whiteList := []string{
		"/admin/open/auth/login",      // 登录接口
		"/admin/open/developer/apply", // 申请开发者资质
	}
	for _, whitePath := range whiteList {
		if strings.HasPrefix(path, whitePath) {
			return true
		}
	}
	return false
}
