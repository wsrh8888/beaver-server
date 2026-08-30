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

import (
	"beaver/common/models"
)

// 存储会话级别的信息（所有用户共享）
// 记录会话类型（私聊/群聊）
// 记录会话的最新消息序列号
// 用于数据同步（客户端需要知道哪些会话有更新）

// ChatConversationMeta 数据同步模型
type ChatConversationMeta struct {
	models.Model
	ConversationID string `gorm:"size:128;uniqueIndex" json:"conversationId"` // 唯一会话ID（私聊/群聊/系统）
	Type           int    `gorm:"not" json:"type"`                            // 1=私聊 2=群聊 3=系统会话
	MaxSeq         int64  `gorm:"not;default:0" json:"maxSeq"`                // 会话全局最新消息序号
	LastMessage    string `gorm:"size:256" json:"lastMessage"`                // 会话最后一条消息预览（全局唯一）
	Version        int64  `gorm:"not;default:0;index" json:"version"`         // 会话元信息版本，用于同步（基于ConversationID递增，从1开始）
}
