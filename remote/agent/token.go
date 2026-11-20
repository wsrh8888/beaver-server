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
	"errors"
	"fmt"
	"net/url"

	"beaver/commonconfig"
	"beaver/utils/ajax"
)

// tokenResp 是 Python 侧 /machine/v1/token 的裸 JSON 回包（非 go-zero 信封）。
// 字段名 snake_case，与对端 OAuth2 风格一致；不要套用 types/remote.TokenResponse（那是信封+lowerCamel）。
type tokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetTokenApi 用 client_credentials 换一次出站令牌，返回 access_token 与有效期（秒，0 表示对端没给）。
//
// beaver-ai 的 token 端点是 Python 实现：snake_case form + 裸 JSON access_token，
// 与 go-zero 机器面（lowerCamel JSON + {code,msg,result} 信封）不同，必须分开处理。
func GetTokenApi(ctx context.Context, clientID, clientSecret string) (string, int, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)

	out, err := ajax.PostForm[tokenResp](ctx, commonconfig.AgentUrl+"/beaver-ai/api/ai/machine/v1/token", nil, form)
	if err != nil {
		return "", 0, fmt.Errorf("换票失败: %w", err)
	}
	if out.AccessToken == "" {
		return "", 0, errors.New("换票响应缺少 access_token")
	}
	return out.AccessToken, out.ExpiresIn, nil
}
