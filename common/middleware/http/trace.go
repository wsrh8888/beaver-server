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

package httpMiddleware

import (
	"net/http"

	"beaver/common/traceid"
)

// TraceMiddleware 使用网关下发的服务端 X-Request-Id；客户端 Uuid 只记日志，不混用。
func TraceMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientUuid := traceid.ClientUuidFromHeaders(r.Header)
		// 只认服务端头；直连调试时没有则本地补一个（不碰 Uuid）
		serverTrace := traceid.Ensure(traceid.ServerTraceFromHeaders(r.Header))
		traceid.AttachServerTrace(w, r, serverTrace)

		ctx := traceid.WithContext(r.Context(), serverTrace, clientUuid)
		next(w, r.WithContext(ctx))
	}
}
