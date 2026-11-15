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

type EmojiPackage struct {
	models.Model
	PackageID   string `gorm:"column:package_id;size:64;uniqueIndex" json:"packageId"` // 全局唯一ID
	Title       string `json:"title"`                                                  // 表情包名称
	CoverFile   string `json:"coverFile"`                                              // 表情包封面文件
	UserID      string `json:"userID"`                                                 // 创建者ID（用户ID或官方ID）
	Description string `json:"description"`                                            // 表情包描述
	Type        string `json:"type"`                                                   // 类型：official-官方，user-用户自定义
	Status      int8   `gorm:"default:1" json:"status"`                                // 状态：1=正常 2=审核中 3=违规禁用
	Version     int64  `gorm:"not null;default:0;index" json:"version"`                // 表情包版本号，每次修改递增
}
