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

import "gorm.io/gorm"

// OpenOAuthToken 用户授权令牌表（含 refresh token）
type OpenOAuthToken struct {
	gorm.Model
	AppID                 string `gorm:"type:varchar(64);index;not null;comment:应用ID"`
	Token                 string `gorm:"type:varchar(256);uniqueIndex;not null;comment:访问令牌"`
	RefreshToken          string `gorm:"type:varchar(256);index;comment:刷新令牌"`
	ExpiresAt             int64  `gorm:"type:bigint;not null;comment:access_token过期时间戳"`
	RefreshTokenExpiresAt int64  `gorm:"type:bigint;not null;comment:refresh_token过期时间戳"`
	Scope                 string `gorm:"type:text;comment:授权范围"`
	UserID                string `gorm:"type:varchar(64);comment:用户ID"`
	OpenID                string `gorm:"type:varchar(64);index;comment:授权用户唯一标识"`
	UnionID               string `gorm:"type:varchar(64);index;comment:用户统一标识"`
}
