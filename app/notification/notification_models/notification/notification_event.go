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

	"gorm.io/datatypes"
)

// NotificationEvent 事件主表：用于幂等和审计，跨服务统一存储事件元数据
type NotificationEvent struct {
	models.Model
	EventID    string         `gorm:"size:64;uniqueIndex;not null" json:"eventId"` // 全局事件ID（雪花/ULID）
	EventType  string         `gorm:"size:32;index;not null" json:"eventType"`     // 事件类型：friend_request/moment_like/...
	Category   string         `gorm:"size:32;index;not null" json:"category"`      // 场景分类：social/system/group/moment等
	Version    int64          `gorm:"not null;default:0;index" json:"version"`     // 全局递增版本号（建议取雪花时间或全局序列），用于增量同步/纠偏
	FromUserID *string        `gorm:"size:64;index" json:"fromUserId,omitempty"`   // 事件触发方
	TargetID   *string        `gorm:"size:64;index" json:"targetId,omitempty"`     // 目标对象（如动态ID、群ID）
	TargetType string         `gorm:"size:32;index" json:"targetType"`             // 目标类型：moment/group/user/message等
	Payload    datatypes.JSON `json:"payload"`                                     // 事件扩展数据，前端可直接渲染
	Priority   int8           `gorm:"not null;default:5;index" json:"priority"`    // 优先级（1最高，9最低），便于队列/推送调度
	Status     int8           `gorm:"not null;default:1;index" json:"status"`      // 1=有效 2=撤回/隐藏 3=失效
	DedupHash  string         `gorm:"size:128;index" json:"dedupHash"`             // 去重哈希，用于点赞合并等
}
