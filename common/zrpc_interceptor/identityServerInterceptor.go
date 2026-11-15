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

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// IdentityServerInterceptor 把入站 metadata 里的身份结论还原成 context 里的 Identity，
// 供 rpc logic 用 common/identity.From 读取。
//
// 只做透传还原，不做验签：验签发生在 HTTP 面的 .api 中间件与边缘网关。
// metadata 里没有身份信息（例如纯内部 rpc→rpc 调用）时保持原 ctx，不覆盖上层结论。
func IdentityServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	first := func(key string) string {
		if v := md.Get(key); len(v) > 0 {
			return v[0]
		}
		return ""
	}

	userID := first("userId")
	clientID := first("clientId")
	kind := first("kind")
	if userID == "" && clientID == "" && kind == "" {
		return handler(ctx, req)
	}

	id := identity.Identity{UserID: userID, ClientID: clientID}
	switch kind {
	case "s2s":
		id.Kind = identity.KindS2S
	case "user":
		id.Kind = identity.KindUser
	default:
		if userID != "" {
			id.Kind = identity.KindUser
		}
	}
	if scopes := md.Get("scopes"); len(scopes) > 0 {
		id.Scopes = scopes
	}

	return handler(identity.With(ctx, id), req)
}
