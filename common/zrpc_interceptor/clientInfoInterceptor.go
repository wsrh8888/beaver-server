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

package zrpc_interceptor

import (
	"context"

	"beaver/common/traceid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ClientInfoInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var clientIP, userID string
	if cl := ctx.Value("clientIP"); cl != nil {
		if s, ok := cl.(string); ok {
			clientIP = s
		}
	}
	if uid := ctx.Value("userId"); uid != nil {
		if s, ok := uid.(string); ok {
			userID = s
		}
	}

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.MD{}
	}
	if clientIP != "" {
		md.Set("clientIP", clientIP)
	}
	if userID != "" {
		md.Set("userId", userID)
	}
	if id := traceid.FromContext(ctx); id != "" {
		md.Set(traceid.MetadataKey, id)
	}

	// 必须保留原 ctx，否则会丢掉 HTTP 侧注入的 trace 等字段
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
