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
	"net/http"

	"beaver/app/open/open_rpc/open"
)

// RequireDeveloperMiddleware 已取消独立开发者表：登录用户即可操作开放平台。
type RequireDeveloperMiddleware struct {
	openRpc open.Open
}

func NewRequireDeveloperMiddleware(openRpc open.Open) *RequireDeveloperMiddleware {
	return &RequireDeveloperMiddleware{openRpc: openRpc}
}

func (m *RequireDeveloperMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Beaver-User-Id")
		if userID == "" {
			http.Error(w, `{"code":401,"msg":"未登录"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
