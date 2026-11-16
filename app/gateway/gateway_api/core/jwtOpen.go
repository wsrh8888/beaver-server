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

package core

import (
	"net/http"

	"beaver/utils/jwts"
)

// verifyOpenTicket 验开放平台开发者票：密钥取 AuthOpen.AccessSecret，
// 受众必须是 beaver:developer，签发行是 beaver-open 的 developer 登录。
// 该票只能进门户面（/beaver-open/...）。
func (p Proxy) verifyOpenTicket(r *http.Request) (*jwts.CustomClaims, error) {
	return verifyTicket(r, p.Config.AuthOpen.AccessSecret, jwts.AudienceDeveloper, "网关未配置开放平台票密钥")
}
