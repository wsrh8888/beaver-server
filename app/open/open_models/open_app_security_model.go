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
	"gorm.io/gorm"
)

// ==================== Security 配置表 ====================

// OpenAppSecurity 安全配置表（对标钉钉开放平台）
type OpenAppSecurity struct {
	gorm.Model
	AppID string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`

	// IP 白名单
	IPWhitelist    string `gorm:"type:text;comment:IP白名单(JSON数组)"`
	TrustedDomains string `gorm:"type:text;comment:可信域名列表(JSON数组)"`

	// 限流配置
	RateLimitEnabled bool `gorm:"type:tinyint;default:0;comment:是否启用限流 1是 0否"`
	RateLimitQPS     int  `gorm:"type:int;default:100;comment:每秒请求数限制"`

	// CSRF 保护
	CSRFProtection bool `gorm:"type:tinyint;default:0;comment:是否启用CSRF保护 1是 0否"`

	// CORS 配置
	AllowedOrigins string `gorm:"type:text;comment:CORS允许的源(JSON数组)"`

	// HTTPS 强制
	RequireHTTPS bool `gorm:"type:tinyint;default:0;comment:是否强制HTTPS 1是 0否"`
}
