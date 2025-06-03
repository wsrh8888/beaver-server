package main

import (
	"beaver/app/gateway/gateway_api/core"
	"beaver/app/gateway/gateway_api/types"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/negroni"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "gateway.yaml", "the config file")

func main() {
	fmt.Println("启动gateway服务")
	flag.Parse()

	var config types.Config

	conf.MustLoad(*configFile, &config)
	logx.SetUp(config.Log)
	proxy := core.Proxy{
		Config: config,
	}

	// 创建带有CORS设置的多路复用器
	n := negroni.New()
	n.Use(negroni.HandlerFunc(corsMiddleware))
	n.UseHandler(proxy)

	err := http.ListenAndServe(config.Addr, n)
	if err != nil {
		fmt.Printf("HTTP服务启动失败: %v\n", err)
		os.Exit(1)
	}
}

func corsMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*") // 设置允许的源域名
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Token, Uuid, Deviceid")

	// 如果是预检请求，直接返回
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 处理实际请求
	next(w, r)
}
