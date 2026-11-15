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

package group_models

import (
	"beaver/common/models"
	"time"
)

// 按照群组独立递增版本
// GroupJoinRequestModel 入群申请表
type GroupJoinRequestModel struct {
	models.Model
	GroupID         string     `gorm:"size:64;index" json:"groupId"`
	ApplicantUserID string     `gorm:"size:64;index" json:"applicantUserId"`
	Message         string     `gorm:"type:text" json:"message"`
	Status          int8       `gorm:"not null;default:0" json:"status"` // 0待审 1同意 2拒绝
	HandledBy       string     `gorm:"size:64" json:"handledBy"`
	HandledAt       *time.Time `json:"handledAt"` // 处理时间
	Version         int64      `gorm:"not null;default:0;index" json:"version"`
}
