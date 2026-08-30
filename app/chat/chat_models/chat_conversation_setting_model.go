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

// 存储用户级别的会话设置（每个用户独立）
// 记录用户对会话的个性化操作（置顶、免打扰、删除等）
// 记录用户在该会话中的已读状态
// 用于UI显示（会话列表、未读消息等）

type ChatUserConversation struct {
	models.Model
	UserID         string `gorm:"size:64;index" json:"userId"`          // 用户ID
	ConversationID string `gorm:"size:128;index" json:"conversationId"` // 关联的会话ID
	IsHidden       bool   `gorm:"default:false" json:"isHidden"`        // 是否在当前用户的会话列表隐藏
	IsPinned       bool   `gorm:"default:false" json:"isPinned"`        // 置顶
	IsMuted        bool   `gorm:"default:false" json:"isMuted"`         // 免打扰
	UserReadSeq    int64  `gorm:"not;default:0" json:"userReadSeq"`     // 当前用户已读游标
	Version        int64  `gorm:"not;default:0;index" json:"version"`   // 用户会话设置版本，用于多端同步（基于UserID递增，所有会话设置共享版本号，从1开始）
}
