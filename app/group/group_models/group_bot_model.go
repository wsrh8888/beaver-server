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

package group_models

import "beaver/common/models"

// GroupBotModel 群内机器人展示模型（只存群内特有信息）
// 基础信息（昵称、头像等）从 user 表获取
// 安全凭证（Token、签名等）在 open_bots 表管理
type GroupBotModel struct {
	models.Model
	GroupID   string `gorm:"size:128;index;not null" json:"groupId"`     // 群组ID
	BotID     string `gorm:"size:128;uniqueIndex;not null" json:"botId"` // 机器人用户ID（关联 users.user_id）
	Status    int    `gorm:"default:1" json:"status"`                    // 1启用 0禁用
	Type      string `gorm:"size:32;default:'custom'" json:"type"`       // 集成类型：custom/github/gitlab/jenkins/grafana/prometheus
	CreatorID string `gorm:"size:128" json:"creatorId"`                  // 创建者用户ID
}
