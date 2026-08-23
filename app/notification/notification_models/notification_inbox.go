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

// NotificationInbox 用户收件箱：一行代表"某用户收到某事件"的状态
type NotificationInbox struct {
	models.Model
	UserID    string     `gorm:"size:64;index:idx_inbox_user_event,priority:1;not null" json:"userId"` // 收件人
	EventID   string     `gorm:"size:64;index:idx_inbox_user_event,priority:2;not null" json:"eventId"`
	EventType string     `gorm:"size:32;index;not null" json:"eventType"` // 冗余字段便于分桶查询
	Category  string     `gorm:"size:32;index;not null" json:"category"`  // 冗余字段便于分桶查询
	Version   int64      `gorm:"not null;default:0;index" json:"version"` // 按用户递增的版本号，用于客户端增量同步
	IsRead    bool       `gorm:"not null;default:false;index" json:"isRead"`
	ReadAt    *time.Time `json:"readAt,omitempty"`
	Status    int8       `gorm:"not null;default:1;index" json:"status"`        // 1=正常 2=隐藏/撤回 3=过期
	IsDeleted bool       `gorm:"not null;default:false;index" json:"isDeleted"` // 用户是否删除该通知
	Silent    bool       `gorm:"not null;default:false" json:"silent"`          // 是否静默（不推送，只计红点）
}
