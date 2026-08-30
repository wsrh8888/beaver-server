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

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// LogMiddleware 自定义的中间件
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ClientIP := httpx.GetRemoteAddr(r)
		ctx := context.WithValue(r.Context(), "ClientIP", ClientIP)

		// 获取请求中的原始域名信息
		originalHost := r.Host
		ctx = context.WithValue(ctx, "ClientHost", originalHost)

		// 判断请求协议是否为 HTTPS 并记录 "http" 或 "https"
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		ctx = context.WithValue(ctx, "Scheme", scheme)

		fmt.Println("scheme:", scheme, "ClientIP:", ClientIP, "ClientHost:", originalHost)
		next(w, r.WithContext(ctx))
	}
}
