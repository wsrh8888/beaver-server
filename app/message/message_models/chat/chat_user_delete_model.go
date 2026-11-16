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

package chat_models

import "beaver/common/models"

// ChatUserDelete 记录用户主动删除的消息（仅自己不可见）
// 对标大厂：不直接修改 ChatMessage 表，而是记录用户的删除行为
type ChatUserDelete struct {
	models.Model
	UserID    string `gorm:"size:64;index:idx_user_msg" json:"userId"`
	MessageID string `gorm:"size:64;index:idx_user_msg" json:"messageId"`
	Version   int64  `gorm:"not null;default:0;index" json:"version"` // 版本号，用于多端同步
}
