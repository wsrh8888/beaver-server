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

package emoji_models

import "beaver/common/models"

// 用户收藏的表情
type EmojiCollectEmoji struct {
	models.Model
	EmojiCollectID string `gorm:"column:emoji_collect_id;size:64;uniqueIndex" json:"emojiCollectId"` // 全局唯一ID
	UserID         string `json:"userId"`                                                            // 用户ID
	EmojiID        string `gorm:"size:64;index" json:"emojiId"`                                      // 表情ID
	IsDeleted      bool   `gorm:"default:false;index" json:"isDeleted"`                              // 是否已删除（软删除）
	Version        int64  `gorm:"not null;default:0;index" json:"version"`                           //基于userId递增
	PackageID      string `gorm:"size:64;index" json:"packageId"`                                    // 表情包ID
}
