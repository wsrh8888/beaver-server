package main

import (
	"beaver/app/gateway/gateway_admin/core"
	"beaver/app/gateway/gateway_admin/types"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/urfave/negroni"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "gateway.yaml", "the config file")

func main() {
	fmt.Println("启动管理后台网关服务 (gateway_admin)")
	flag.Parse()

	var config types.Config
	conf.MustLoad(*configFile, &config)
	logx.SetUp(config.Log)

	// 初始化代理
	proxy := &core.Proxy{
		Config: config,
	}

	// 初始化反向代理
	if err := proxy.Init(); err != nil {
		fmt.Printf("初始化网关失败: %v\n", err)
		os.Exit(1)
	}

	// 创建带有CORS设置的多路复用器
	n := negroni.New()
	n.Use(negroni.HandlerFunc(corsMiddleware))
	n.Use(negroni.HandlerFunc(requestLogMiddleware))
	n.UseHandler(proxy)

	fmt.Printf("网关服务启动在 %s\n", config.Addr)
	if config.Etcd != "" {
		fmt.Printf("后端服务发现（etcd）: %s, 服务key: backend_admin\n", config.Etcd)
	} else {
		fmt.Printf("警告: 未配置 etcd\n")
	}
	if config.Auth.AccessSecret != "" {
		fmt.Printf("JWT 认证已启用（网关直接解析 JWT，提升性能）\n")
	}
	if config.Redis.Addr != "" {
		fmt.Printf("Redis 验证已启用（将验证 token 在 Redis 中的有效性）\n")
	}
	err := http.ListenAndServe(config.Addr, n)
	if err != nil {
		fmt.Printf("HTTP服务启动失败: %v\n", err)
		os.Exit(1)
	}
}

// requestLogMiddleware 请求日志中间件
func requestLogMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(w, r)
	duration := time.Since(start)
	logx.Infof("[%s] %s %s - %v", r.Header.Get("Uuid"), r.Method, r.URL.Path, duration)
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
