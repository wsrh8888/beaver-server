package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"beaver/app/doc/doc_api/internal/config"
	"beaver/app/doc/doc_api/internal/handler"
	"beaver/app/doc/doc_api/internal/svc"
	"beaver/common/etcd"
	"beaver/common/middleware"
	httpMiddleware "beaver/common/middleware/http"
	"beaver/utils/beaverlog"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/doc.yaml", "the config file")

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
		fmt.Printf("doc_api 协同网关关闭异常: %v\n", err)
	}
	server.Stop()
}
