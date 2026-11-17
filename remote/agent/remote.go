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

package agent

import (
	"context"
	"fmt"
	"net/http"

	"beaver/utils/ajax"
)

// 本包唯一的出站通道：静态 X-Internal-Token 鉴权 + 带令牌发包，业务方法在各自文件里。
// 与 open/ai 等 peer 同构：走 ajax.Post 统一按 go-zero 信封 {code,msg,result} 判成败、
// 把 result 解成 Res；本包只负责补 X-Internal-Token 头。

// Post 发一次 POST 的 JSON 请求，补静态 X-Internal-Token，响应按信封 result 解成 Res。
func Post[Res, Req any](ctx context.Context, endpoint string, header http.Header, body Req) (*Res, error) {
	h, err := authHeader(ctx, header)
	if err != nil {
		return nil, fmt.Errorf("%s：装取出站令牌失败: %w", endpoint, err)
	}
	return ajax.Post[Res](ctx, endpoint, "", h, body)
}

// authHeader 把调用方给的头复制一份，再补上 agent 的静态 InternalToken。
func authHeader(ctx context.Context, header http.Header) (http.Header, error) {
	h := make(http.Header, len(header)+1)
	for k, vs := range header {
		for _, v := range vs {
			h.Add(k, v)
		}
	}
	if client != nil && client.token != "" {
		h.Set("X-Internal-Token", client.token)
	}
	return h, nil
}
