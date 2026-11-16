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

package notification_models

import (
	"beaver/common/models"
	"time"
)

// NotificationRead 用户分类已读游标：记录用户对每个分类的最后查看时间
// 用于支持增量同步和状态管理
type NotificationRead struct {
	models.Model
	UserID     string     `gorm:"size:64;uniqueIndex:uniq_cursor;index:idx_cursor_user_category,priority:1;not null" json:"userId"`
	Category   string     `gorm:"size:32;uniqueIndex:uniq_cursor;index:idx_cursor_user_category,priority:2;not null" json:"category"`
	Version    int64      `gorm:"not null;default:0;index" json:"version"` // 游标版本（按用户+分类递增），便于幂等/同步
	LastReadAt *time.Time `json:"lastReadAt,omitempty"`                    // 最后查看该分类的时间
}
