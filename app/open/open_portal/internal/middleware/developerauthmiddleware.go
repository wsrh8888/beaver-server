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
	"net/http"
	"strings"

	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeveloperAuthMiddleware struct {
	secretKey string
}

func NewDeveloperAuthMiddleware(secretKey string) *DeveloperAuthMiddleware {
	return &DeveloperAuthMiddleware{
		secretKey: secretKey,
	}
}

func (m *DeveloperAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromRequest(r)
		if token == "" {
			http.Error(w, `{"code":401,"msg":"缺少认证令牌"}`, http.StatusUnauthorized)
			return
		}

		claims, err := jwts.ParseToken(token, m.secretKey)
		if err != nil {
			logx.Errorf("JWT 验证失败: %v", err)
			http.Error(w, `{"code":401,"msg":"令牌无效或已过期"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		r = r.WithContext(ctx)
		r.Header.Set("Beaver-User-Id", claims.UserID)
		next(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	if token := r.Header.Get("Token"); token != "" {
		return token
	}

	return r.URL.Query().Get("token")
}
