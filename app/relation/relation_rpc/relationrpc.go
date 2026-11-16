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

	"beaver/app/relation/relation_rpc/internal/config"
	friendServer "beaver/app/relation/relation_rpc/internal/server/friend"
	groupServer "beaver/app/relation/relation_rpc/internal/server/group"
	"beaver/app/relation/relation_rpc/internal/svc"
	"beaver/app/relation/relation_rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/relationrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		// goctl multiple 模式：一个 proto（relation.proto）、一个进程、一个 etcd 注册，
		// 同一个 grpc server 挂载 friend + group 两套 service（server/logic/client 按子域分目录）。
		pb.RegisterFriendServer(grpcServer, friendServer.NewFriendServer(ctx))
		pb.RegisterGroupServer(grpcServer, groupServer.NewGroupServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 还原 HTTP 面验签后的身份结论，rpc logic 只从 context 取。
	s.AddUnaryInterceptors(zrpc_interceptor.IdentityServerInterceptor)

	fmt.Printf("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}
