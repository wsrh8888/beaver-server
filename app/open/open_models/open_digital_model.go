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
	DigitalRuntimeBuiltin   = "builtin"   // 提示词模式：内置 IM + beaver-agent 跑 LLM，编排配置在 OpenDigitalConfig
	DigitalRuntimeExternal = "external"  // 长链接模式：对接外部 agent（龙虾等），走 WebSocket 转发
	DigitalRuntimeWorkflow  = "workflow" // 工作流模式：按预设节点流程编排，多步骤/多分支/条件流转
	// DigitalRuntimePlugin = "plugin" // 插件模式（预留）
)

const (
	DigitalVisibilityPrivate = "private" // 仅管理员可搜到
	DigitalVisibilityPartial = "partial" // 指定人可搜到（role=admin|viewer）
	DigitalVisibilityPublic  = "public"  // 全 IM 用户可搜到
)

const (
	DigitalStatusEnabled  = 1
	DigitalStatusDisabled = 0
)

const (
	DigitalSourceSystem = "system" // 管理后台创建的系统数字伙伴
	DigitalSourceOpen   = "open"   // 开放平台创建
)

// OpenDigital 开放平台三大能力之一：数字伙伴（一级实体，不挂 OpenOAuth）。
//
// DigitalID = UserModel.UserID；名字/头像/简介只在 UserModel，本表不存展示字段。
// external：AppKey + AppSecret（长期凭据，不过期，仅手动重置）；builtin：AppKey/密钥为空串。
// 调 Open API / 建连：用 AppKey+AppSecret 换短时 access_token（见 OpenDigitalToken）。
type OpenDigital struct {
	gorm.Model
	DigitalID string `gorm:"column:digital_id;type:varchar(64);uniqueIndex;not null;comment:数字伙伴 UserModel.UserID"`
	OwnerID   string `gorm:"column:owner_id;type:varchar(64);index;not null;comment:创建者用户ID"`

	Runtime string `gorm:"type:varchar(16);not null;comment:builtin|external"`

	// AppKey 仅 external 非空，服务端生成（dig_<uuid>，唯一性由生成保证）；builtin 为空串，故不建 uniqueIndex。
	AppKey    string `gorm:"column:app_key;type:varchar(64);index;not null;default:'';comment:外接 AppKey dig_xxx"`
	AppSecret string `gorm:"column:app_secret_hash;type:varchar(128);not null;default:'';comment:外接 AppSecret 哈希"`

	Visibility string `gorm:"type:varchar(16);not null;index;comment:private|partial|public"`
	Status     int    `gorm:"column:status;type:tinyint;not null;default:1;index;comment:1启用 0禁用"`
	Source     string `gorm:"column:source;type:varchar(16);not null;default:'open';index;comment:system系统创建 open开放平台创建"`
}

func (OpenDigital) TableName() string { return "open_digital" }
