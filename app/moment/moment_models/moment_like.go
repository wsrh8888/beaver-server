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

package moment_models

import (
	"beaver/common/models"
)

/**
 * @description: 点赞表
 */
type MomentLikeModel struct {
	models.Model
	LikeID    string `gorm:"column:like_id;size:64;uniqueIndex;not null" json:"likeId"` // 全局唯一ID
	MomentID  string `gorm:"size:64;not null;index" json:"momentId"`                    // 动态ID (关联 moment_id)
	UserID    string `gorm:"size:64;not null;index" json:"userId"`                      // 点赞用户Id (索引)
	IsDeleted bool   `gorm:"not null;default:false;index" json:"isDeleted"`             // 软删除标记 (索引)
}
