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

package authgate

import (
	"encoding/json"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)

// unauthorizedBody 与 common/response.Body 同形，便于调用方统一解包。
type unauthorizedBody struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

// New 把闸门接成 go-zero 的 rest.Middleware。
//
// 用法：由 .api 的 `@server middleware: AuthGate` 声明，goctl 生成
// `rest.WithMiddlewares([]rest.Middleware{serverCtx.AuthGate}, ...)`；
// 各服务 ServiceContext 用本函数装配同名字段。禁止再用 main 的 server.Use 全局挂载。
func New(g *Gate) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			id, ok := g.Resolve(r)
			if !ok {
				writeUnauthorized(w, "缺少或无效的访问凭据")
				return
			}
			next(w, g.Apply(r, id))
		}
	}
}

// writeUnauthorized 输出真实 HTTP 401。
//
// 不走 common/response.Response：那个封装恒发 200 + code=1，是业务错误语义；
// 鉴权失败属于传输层结论，S2S 客户端要靠 401 触发清票重换。
func writeUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(unauthorizedBody{Code: 1, Msg: msg, Result: nil})
}
