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

	"beaver/app/ws/ws_api/internal/config"
	"beaver/middleware/authgate"

	"github.com/zeromicro/go-zero/rest"
)

// AuthGateMiddleware 由 .api 的 `@server middleware: AuthGate` 声明接线。
// 实现共享在 beaver/middleware/authgate，这里只负责用本服务配置装配凭据。
type AuthGateMiddleware struct {
	gate rest.Middleware
}

func NewAuthGateMiddleware(c config.Config) *AuthGateMiddleware {
	g := authgate.NewGate(c.Auth.AccessSecret, c.S2S.Secret, authgate.AudienceCore, authgate.WithQueryToken())
	if err := g.Validate(); err != nil {
		panic("ws_api 身份闸门配置不完整: " + err.Error())
	}
	return &AuthGateMiddleware{gate: authgate.New(g)}
}

func (m *AuthGateMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.gate(next)
}
