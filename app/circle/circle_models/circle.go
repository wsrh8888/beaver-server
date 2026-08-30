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

package circle_models

import "beaver/common/models"

// CircleModel 圈子主表
type CircleModel struct {
	models.Model
	CircleID    string `gorm:"column:circle_id;size:64;uniqueIndex;not null" json:"circleId"` // 圈子唯一ID
	Name        string `gorm:"size:64;not null" json:"name"`                                  // 圈子名称
	Description string `gorm:"type:text" json:"description"`                                  // 圈子简介
	Avatar      string `gorm:"size:256" json:"avatar"`                                        // 圈子头像
	CreatorID   string `gorm:"size:64;not null;index" json:"creatorId"`                       // 创建者用户ID
	JoinType    int8   `gorm:"not null;default:0" json:"joinType"`                            // 加入方式：0=自由加入 1=审批加入
	Version     int64  `gorm:"not null;default:0;index" json:"version"`                       // 版本号，用于客户端增量同步
	IsDeleted   bool   `gorm:"not null;default:false;index" json:"isDeleted"`                 // 软删除标记
}
