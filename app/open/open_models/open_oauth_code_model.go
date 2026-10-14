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

package open_models

import (
	"time"

	"gorm.io/gorm"
)

// OpenOAuthCode OAuth 授权码（短生命周期、一次性；生产可迁 Redis，本表作持久/审计兜底）。
// Code 仅存哈希；换票时校验 code_challenge / method。
type OpenOAuthCode struct {
	gorm.Model
	Code                string     `gorm:"column:code_hash;type:varchar(128);uniqueIndex;not null;comment:授权码哈希"`
	OAuthID             string     `gorm:"column:oauth_id;type:varchar(64);index;not null;comment:OpenOAuth.OAuthID"`
	UserID              string     `gorm:"column:user_id;type:varchar(64);index;not null;comment:用户ID"`
	RedirectURI         string     `gorm:"column:redirect_uri;type:varchar(512);comment:回调地址"`
	Scope               string     `gorm:"type:varchar(512);comment:权限范围"`
	State               string     `gorm:"type:varchar(128);comment:CSRF state"`
	CodeChallenge       string     `gorm:"column:code_challenge;type:varchar(128);comment:PKCE code_challenge"`
	CodeChallengeMethod string     `gorm:"column:code_challenge_method;type:varchar(16);comment:PKCE method S256|plain"`
	Nonce               string     `gorm:"type:varchar(128);comment:OIDC nonce"`
	ExpiresAt           time.Time  `gorm:"column:expires_at;type:datetime;index;not null;comment:过期时间"`
	UsedAt              *time.Time `gorm:"column:used_at;type:datetime;index;comment:使用时间，空表示未使用"`
	Scene               string     `gorm:"type:varchar(20);not null;default:oauth;index;comment:场景 oauth/h5_sso/pc_scan"`
}

func (OpenOAuthCode) TableName() string { return "open_oauth_codes" }
