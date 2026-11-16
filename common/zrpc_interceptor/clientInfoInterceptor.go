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

	"beaver/common/identity"
	"beaver/common/traceid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ClientInfoInterceptor 把 HTTP 面的调用信息（客户端 IP、已确认的身份结论）透传到 gRPC metadata。
//
// 身份不由这里生成：HTTP 面只挂 JwtS2s 验机器票，用户身份采信网关写入的 Beaver-User-Id /
// Beaver-Admin-Id 头，由 .api 的 Req 字段（header tag）直接绑定，logic 之间以显式参数传递。
// 因此 identity.From 通常为空；若上游确实用 identity.With 写过结论，这里原样透传，
// rpc 侧由 IdentityServerInterceptor 还原成同一份 Identity。
func ClientInfoInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var clientIP string
	if cl := ctx.Value("clientIP"); cl != nil {
		if s, ok := cl.(string); ok {
			clientIP = s
		}
	}

	id := identity.From(ctx)

	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.MD{}
	}
	if clientIP != "" {
		md.Set("clientIP", clientIP)
	}
	if id.UserID != "" {
		md.Set("userId", id.UserID)
	}
	if id.ClientID != "" {
		md.Set("clientId", id.ClientID)
	}
	switch id.Kind {
	case identity.KindUser:
		md.Set("kind", "user")
	case identity.KindS2S:
		md.Set("kind", "s2s")
	}
	if len(id.Scopes) > 0 {
		md.Set("scopes", id.Scopes...)
	}
	if id := traceid.FromContext(ctx); id != "" {
		md.Set(traceid.MetadataKey, id)
	}

	// 必须保留原 ctx，否则会丢掉 HTTP 侧注入的 trace 等字段
	ctx = metadata.NewOutgoingContext(ctx, md)
	return invoker(ctx, method, req, reply, cc, opts...)
}
