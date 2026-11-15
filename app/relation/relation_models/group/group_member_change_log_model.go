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

// 全局递增版本
// 群成员变更日志
type GroupMemberChangeLogModel struct {
	models.Model
	GroupID    string    `gorm:"size:64;index" json:"groupId"`  // 群ID
	UserID     string    `gorm:"size:64;index" json:"userId"`   // 用户ID
	ChangeType string    `gorm:"size:32" json:"changeType"`     // join/leave/kick/promote/demote
	OperatedBy string    `gorm:"size:64" json:"operatedBy"`     // 操作者（群主/管理员）
	ChangeTime time.Time `json:"changeTime"`                    // 变更时间
	Version    int64     `gorm:"not null;index" json:"version"` // 版本号
}
