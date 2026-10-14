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

const (
	DigitalMemberRoleAdmin  = "admin"  // 管理员：可管理 + 可搜索
	DigitalMemberRoleViewer = "viewer" // 可见者：仅可搜索
)

// OpenDigitalMember 数字伙伴成员（管理员 + 可见者合并一表，用 Role 区分）。
// Owner 天然有管理权，落主表 OwnerID，不入本表。
// Visibility=private 时仅 Owner + admin 可搜；partial 时 admin+viewer 可搜；public 全员可搜。
type OpenDigitalMember struct {
	gorm.Model
	DigitalID string `gorm:"column:digital_id;type:varchar(64);uniqueIndex:uk_digital_member;not null;comment:数字伙伴ID"`
	UserID    string `gorm:"column:user_id;type:varchar(64);uniqueIndex:uk_digital_member;index;not null;comment:成员用户ID"`
	Role      string `gorm:"column:role;type:varchar(16);not null;default:viewer;comment:角色 admin|viewer"`
}

func (OpenDigitalMember) TableName() string { return "open_digital_member" }
