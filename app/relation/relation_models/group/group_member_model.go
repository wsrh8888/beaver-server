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

import (
	"beaver/common/models"
	"time"
)

// 群成员（版本号按群组独立递增）
type GroupMemberModel struct {
	models.Model
	GroupID    string     `gorm:"size:64;index" json:"groupId"`
	UserID     string     `gorm:"size:64;index" json:"userId"`
	Role       int8       `json:"role"`                                    // 1群主 2管理员 3普通成员
	Status     int8       `gorm:"default:1" json:"status"`                 // 1正常 2退出 3被踢
	JoinTime   time.Time  `json:"joinTime"`                                // 加入时间
	MutedUntil *time.Time `json:"mutedUntil"`                              // 禁言截止时间，nil表示未禁言
	Version    int64      `gorm:"not null;default:0;index" json:"version"` // 群组成员列表版本号
}
