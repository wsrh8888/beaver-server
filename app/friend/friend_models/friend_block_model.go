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

package friend_models

import "beaver/common/models"

// FriendBlockModel 黑名单表
type FriendBlockModel struct {
	models.Model
	BlockID       string `gorm:"column:block_id;size:64;uniqueIndex;not null" json:"blockId"`
	UserID        string `gorm:"size:64;index;not null" json:"userId"`        // 执行拉黑的用户ID
	BlockedUserID string `gorm:"size:64;index;not null" json:"blockedUserId"` // 被拉黑的用户ID
}
