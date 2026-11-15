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

package auth_models

import (
	"beaver/common/models"
	"time"
)

// AuthCredentialModel 认证凭证模型（密码、登录记录等敏感信息）
type AuthCredentialModel struct {
	models.Model
	UserID      string     `gorm:"size:64;uniqueIndex;not null" json:"userId"` // 关联 users.user_id
	Password    string     `gorm:"size:128;not null" json:"-"`                 // 密码哈希（不对外暴露）
	Salt        string     `gorm:"size:32" json:"-"`                           // 盐值（可选，如果用 bcrypt 则不需要）
	LastLoginAt *time.Time `json:"lastLoginAt"`                                // 最后登录时间
	LoginCount  int64      `gorm:"default:0" json:"loginCount"`                // 登录次数
}
