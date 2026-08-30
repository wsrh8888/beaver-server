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

	"beaver/app/open/open_rpc/open"
	"beaver/app/open/open_rpc/types/open_rpc"
)

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

		res, err := m.openRpc.GetDeveloperByUserID(r.Context(), &open_rpc.GetDeveloperByUserIDReq{UserId: userID})
		if err != nil || !res.Found || res.Developer == nil || res.Developer.Status != 1 {
			http.Error(w, `{"code":403,"msg":"您还不是认证开发者,请先申请开发者资质"}`, http.StatusForbidden)
			return
		}

		next(w, r.WithContext(context.WithValue(r.Context(), "developerId", res.Developer.Id)))
	}
}
