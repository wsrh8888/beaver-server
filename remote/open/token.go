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
	"errors"
	"fmt"
	"net/url"

	"beaver/commonconfig"
	remotetype "beaver/types/remote"
	"beaver/utils/ajax"
)

// GetToken 用 client_credentials 换一次出站令牌，返回 access_token 与有效期（秒，0 表示对端没给）。
//
// 换票是 beaver-open 的一个独立接口，地址和表单都在这里，调用方只管拿票。
func GetToken(ctx context.Context, clientID, clientSecret string) (string, int, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	out, err := ajax.PostForm[remotetype.TokenResponse](ctx, commonconfig.OpenUrl+"/internal/token", nil, form)
	if err != nil {
		return "", 0, fmt.Errorf("换票失败: %w", err)
	}
	if out.AccessToken == "" {
		return "", 0, errors.New("换票响应缺少 access_token")
	}
	return out.AccessToken, out.ExpiresIn, nil
}
