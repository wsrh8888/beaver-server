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
	"beaver/common/zrpc_interceptor"
	"flag"
	"fmt"

	"beaver/app/message/message_rpc/internal/config"
	"beaver/app/message/message_rpc/internal/logic/consumer"
	conversationServer "beaver/app/message/message_rpc/internal/server/conversation"
	messageServer "beaver/app/message/message_rpc/internal/server/message"
	"beaver/app/message/message_rpc/internal/svc"
	chat_rpc "beaver/app/message/message_rpc/pb"
	grpcMiddleware "beaver/common/middleware/grpc"
	"beaver/utils/beaverlog"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/messagerpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	beaverlog.InitFromConf(c.RpcServerConf.ServiceConf)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		chat_rpc.RegisterConversationServer(grpcServer, conversationServer.NewConversationServer(ctx))
		chat_rpc.RegisterMessageServer(grpcServer, messageServer.NewMessageServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 添加拦截器
	s.AddUnaryInterceptors(grpcMiddleware.RequestLogInterceptor, zrpc_interceptor.IdentityServerInterceptor)

	// 启动 RocketMQ Consumer（agent_task 编排层 + agent_done 终稿落库）
	if ctx.RocketMQ != nil {
		agentTaskConsumer := consumer.NewAgentTaskConsumerLogic(ctx)
		go func() {
			if err := agentTaskConsumer.StartConsumer(); err != nil {
				fmt.Printf("agent_task Consumer 启动失败: %v\n", err)
			}
		}()
		agentDoneConsumer := consumer.NewAgentDoneConsumerLogic(ctx)
		go func() {
			if err := agentDoneConsumer.StartConsumer(); err != nil {
				fmt.Printf("agent_done Consumer 启动失败: %v\n", err)
			}
		}()
	}

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
