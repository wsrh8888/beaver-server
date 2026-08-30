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
	"flag"
	"fmt"

	"beaver/app/ws/ws_api/internal/config"
	"beaver/app/ws/ws_api/internal/handler"
	"beaver/app/ws/ws_api/internal/logic"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/common/etcd"
	"beaver/common/middleware"
	"beaver/utils/logger"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ws.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logger.Init("ws_api")

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 启动 RocketMQ Consumer
	if ctx.RocketMQ != nil {
		mqConsumer := logic.NewMqConsumerLogic(ctx)
		go func() {
			if err := mqConsumer.StartConsumer(); err != nil {
				fmt.Printf("RocketMQ Consumer 启动失败: %v\n", err)
			}
		}()
	}

	etcd.DeliveryAddress(c.Etcd, c.Name+"_api", fmt.Sprintf("%s:%d", c.Host, c.Port))

	server.Use(middleware.LogMiddleware)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
