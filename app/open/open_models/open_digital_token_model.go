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

// OpenDigitalToken 外接数字伙伴短时 access_token（Open API / 长连接 ticket）。
// 用 OpenDigital.AppKey + AppSecret 换取；无 refresh，过期后重新换票。
// 关联用一级主键 DigitalID（唯一），不依赖 AppKey；令牌仅存哈希。
type OpenDigitalToken struct {
	gorm.Model
	DigitalID string     `gorm:"column:digital_id;type:varchar(64);index;not null;comment:OpenDigital.DigitalID"`
	Token     string     `gorm:"column:access_token_hash;type:varchar(128);uniqueIndex;not null;comment:access_token 哈希"`
	ExpiresAt time.Time  `gorm:"column:expires_at;type:datetime;index;not null;comment:access_token 过期时间"`
	RevokedAt *time.Time `gorm:"column:revoked_at;type:datetime;index;comment:撤销时间"`
}

func (OpenDigitalToken) TableName() string { return "open_digital_tokens" }
