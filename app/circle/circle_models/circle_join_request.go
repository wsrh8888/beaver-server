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

// CircleJoinRequestModel 加圈申请表
type CircleJoinRequestModel struct {
	models.Model
	CircleID string `gorm:"size:64;not null;index" json:"circleId"`                           // 圈子ID
	UserID   string `gorm:"size:64;not null;index" json:"userId"`                             // 申请用户ID
	Status   int8   `gorm:"not null;default:0" json:"status"`                                 // 状态：0=待审批 1=已通过 2=已拒绝
	Reason   string `gorm:"size:256" json:"reason"`                                           // 申请理由
}
