package main

import (
	"beaver/app/doc/doc_sync_ws/internal/config"
	"beaver/app/doc/doc_sync_ws/internal/handler"
	"beaver/app/doc/doc_sync_ws/internal/svc"
	"beaver/common/etcd"
	"beaver/common/middleware"
	httpMiddleware "beaver/common/middleware/http"
	"beaver/utils/beaverlog"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/docsync.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	beaverlog.InitFromConf(c.RestConf.ServiceConf)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	server.Use(httpMiddleware.TraceMiddleware)
	server.Use(middleware.LogMiddleware)

	go func() {
		fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
		server.Start()
	}()

	// 优雅关闭：收到退出信号后关闭协同网关，卸载房间并触发末次压实落库。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := ctx.Gateway.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("doc_sync_ws 网关关闭异常: %v\n", err)
	}
	server.Stop()
}
