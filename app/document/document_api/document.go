package main

import (
	"flag"
	"fmt"

	"beaver/app/document/document_api/internal/config"
	"beaver/app/document/document_api/internal/handler"
	"beaver/app/document/document_api/internal/svc"
	"beaver/common/etcd"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/document.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	fmt.Printf("Starting document server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
