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

// OpenOAuthToken OAuth 用户授权令牌（authorization_code / 扫码等），仅服务 OpenOAuth。
// Access/Refresh 仅存哈希；身份用 UserID。
// TokenFamilyID 用于 refresh 轮换与 reuse detection。
type OpenOAuthToken struct {
	gorm.Model
	OAuthID               string     `gorm:"column:oauth_id;type:varchar(64);index;not null;comment:OpenOAuth.OAuthID"`
	UserID                string     `gorm:"column:user_id;type:varchar(64);index;not null;comment:授权用户ID"`
	Token                 string     `gorm:"column:access_token_hash;type:varchar(128);uniqueIndex;not null;comment:access_token 哈希"`
	RefreshToken          string     `gorm:"column:refresh_token_hash;type:varchar(128);index;comment:refresh_token 哈希"`
	TokenFamilyID         string     `gorm:"column:token_family_id;type:varchar(64);index;not null;default:'';comment:refresh 令牌族ID"`
	ExpiresAt             time.Time  `gorm:"column:expires_at;type:datetime;index;not null;comment:access_token 过期时间"`
	RefreshTokenExpiresAt time.Time `gorm:"column:refresh_token_expires_at;type:datetime;index;comment:refresh_token 过期时间（无 refresh 时为零值）"`
	RevokedAt             *time.Time `gorm:"column:revoked_at;type:datetime;index;comment:撤销时间"`
	Scope                 string     `gorm:"type:varchar(512);comment:授权范围"`
}

func (OpenOAuthToken) TableName() string { return "open_oauth_tokens" }
