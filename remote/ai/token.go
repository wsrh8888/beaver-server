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

package ai

import (
	"context"
	"errors"
	"fmt"

	"beaver/commonconfig"
	remotetype "beaver/types/remote"
	"beaver/utils/ajax"
)

// GetToken 用 client_credentials 换一次出站令牌，返回 accessToken 与有效期（秒，0 表示对端没给）。
//
// 换票是 beaver-ai / ai 的独立机器视图接口，地址与请求体都在这里，调用方只管拿票。
func GetToken(ctx context.Context, clientID, clientSecret string) (string, int, error) {
	out, err := ajax.Post[remotetype.TokenResponse](ctx, commonconfig.AiUrl, "/beaver-ai/api/ai/machine/v1/token", nil, remotetype.TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return "", 0, fmt.Errorf("换票失败: %w", err)
	}
	if out.AccessToken == "" {
		return "", 0, errors.New("换票响应缺少 accessToken")
	}
	return out.AccessToken, out.ExpiresIn, nil
}
