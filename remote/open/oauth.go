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

package open

import (
	"context"
	"net/http"

	"beaver/commonconfig"
	opentype "beaver/types/remote/open"
)

// ExchangeToken OAuth 授权码换 token。
func ExchangeToken(ctx context.Context, header http.Header, body *opentype.ExchangeTokenReq) (*opentype.ExchangeTokenRes, error) {
	return Post[opentype.ExchangeTokenRes](ctx, commonconfig.OpenUrl+"/internal/Open.ExchangeToken", header, body)
}

// ListOpenOAuth 列开放平台 OAuth 应用。
func ListOpenOAuth(ctx context.Context, header http.Header, body *opentype.ListOpenOAuthReq) (*opentype.ListOpenOAuthRes, error) {
	return Post[opentype.ListOpenOAuthRes](ctx, commonconfig.OpenUrl+"/internal/Open.ListOpenOAuth", header, body)
}

// UpdateOpenOAuth 更新开放平台 OAuth 应用。
func UpdateOpenOAuth(ctx context.Context, header http.Header, body *opentype.UpdateOpenOAuthReq) (*opentype.UpdateOpenOAuthRes, error) {
	return Post[opentype.UpdateOpenOAuthRes](ctx, commonconfig.OpenUrl+"/internal/Open.UpdateOpenOAuth", header, body)
}
