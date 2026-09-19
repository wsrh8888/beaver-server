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

// ==================== OpenApp 应用主表 ====================

// OpenApp 开放平台应用主表
//
// 职责边界：本表**只是接入凭据** —— 它是谁、叫什么、归谁、开没开。
// 「这个应用能做什么」由各能力表表达，不要塞回本表：
//   OAuth 配置   → OpenOAuthConfig
//   智能机器人   → OpenRobot
//   通知机器人   → OpenBotModel
//   （客户端微应用 / JSSDK 鉴权如需支持，应独立成表）
type OpenApp struct {
	gorm.Model
	// 身份认证
	AppID     string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用唯一标识"`
	AppSecret string `gorm:"type:varchar(128);not null;comment:应用密钥"`
	// 基础信息
	Name        string `gorm:"type:varchar(100);not null;comment:应用名称"`
	Description string `gorm:"type:text;comment:应用描述"`
	Icon        string `gorm:"type:varchar(500);comment:应用图标URL"`
	OwnerID     string `gorm:"type:varchar(64);index;comment:所属用户ID"`
	// 状态：只区分「能不能用」。创建即启用，不需要草稿/发布/审核这类流程态。
	Status int `gorm:"type:tinyint;default:1;comment:1启用 0禁用"`
}
