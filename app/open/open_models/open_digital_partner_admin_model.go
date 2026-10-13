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

// OpenDigitalPartnerAdmin 数字伙伴管理员（Owner 天然有管理权，额外管理员落本表）。
type OpenDigitalPartnerAdmin struct {
	gorm.Model
	PartnerID string `gorm:"column:partner_id;type:varchar(64);uniqueIndex:uk_partner_admin;not null;comment:数字伙伴ID"`
	UserID    string `gorm:"column:user_id;type:varchar(64);uniqueIndex:uk_partner_admin;index;not null;comment:管理员用户ID"`
}

func (OpenDigitalPartnerAdmin) TableName() string { return "open_digital_partner_admin" }
