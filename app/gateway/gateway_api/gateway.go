/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package main

import (
	"beaver/app/gateway/gateway_api/core"
	"beaver/app/gateway/gateway_api/types"
	"beaver/utils/logger"
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
	logger.Init("gateway_api")
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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	// 如果是预检请求，直接返回
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 处理实际请求
	next(w, r)
}
