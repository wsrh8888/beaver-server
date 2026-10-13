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
	PartnerRuntimeBuiltin  = "builtin"
	PartnerRuntimeExternal = "external"
)

// 可见范围（搜索发现，无草稿态）
const (
	PartnerVisibilityPrivate = "private" // 非公开：仅管理员可搜到
	PartnerVisibilityPartial = "partial" // 部分公开：指定人可搜到（见 OpenDigitalPartnerViewer）
	PartnerVisibilityPublic  = "public"  // 公开：全 IM 用户可搜到
)

// OpenDigitalPartner 数字伙伴（用户创建，一人可多个）。
//
// PartnerID = 申请到的 UserModel.UserID（与普通用户同一套账号体系，落 user 库）。
// 名字/头像/简介 → UserModel（NickName/Avatar/Abstract），本表不冗余。
type OpenDigitalPartner struct {
	gorm.Model
	// 创建时走 user 发放 UserID，写入本字段；即「工号/账号」本身
	PartnerID   string `gorm:"column:partner_id;type:varchar(64);uniqueIndex;not null;comment:数字伙伴 UserModel.UserID"`
	OwnerID     string `gorm:"column:owner_id;type:varchar(64);index;not null;comment:创建者用户ID"`
	Description string `gorm:"type:varchar(512);not null;default:'';comment:开放平台侧补充描述"`

	Runtime string `gorm:"type:varchar(16);not null;index;comment:builtin|external"`
	// external 必填 OpenApp.AppID；builtin 固定空串
	AppID string `gorm:"column:app_id;type:varchar(64);index;not null;default:'';comment:外接应用ID"`

	Visibility string `gorm:"type:varchar(16);not null;index;comment:private|partial|public"`
	Status     int8   `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用"`
}

func (OpenDigitalPartner) TableName() string { return "open_digital_partner" }
