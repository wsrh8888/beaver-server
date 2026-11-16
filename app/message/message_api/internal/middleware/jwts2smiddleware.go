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

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package middleware

import (
	"net/http"

	"beaver/app/message/message_api/internal/config"
	corejwt "beaver/middleware/jwt"
	"beaver/utils/s2s"
	"github.com/zeromicro/go-zero/rest"
)

// JwtS2sMiddleware 由 .api 的 `@server middleware: JwtS2s` 声明接线。
// 验票中间件在 middleware/jwt（内部走 utils/s2s.Verify），audience 与主仓业务面一致（core:api），这里只负责装配。
type JwtS2sMiddleware struct {
	mw rest.Middleware
}

func NewJwtS2sMiddleware(c config.Config) *JwtS2sMiddleware {
	return &JwtS2sMiddleware{mw: corejwt.NewJwtS2s(c.S2S.Secret, s2s.AudienceCore)}
}

func (m *JwtS2sMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return m.mw(next)
}
