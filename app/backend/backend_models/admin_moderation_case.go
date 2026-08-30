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

package backend_models

import (
	"beaver/common/models"
	"time"
)

// 工单来源
const (
	CaseSourceReport  = 1
	CaseSourceFeedback = 2
	CaseSourceManual  = 3
)

// 工单状态
const (
	CaseStatusPending    = 1
	CaseStatusProcessing = 2
	CaseStatusResolved   = 3
	CaseStatusRejected   = 4
)

// 处置对象类型
const (
	CaseTargetUser    = 1
	CaseTargetMessage = 2
	CaseTargetMoment  = 3
	CaseTargetGroup   = 4
)

// AdminModerationCase 运营处置工单（后台域，跨 RPC 编排入口）
type AdminModerationCase struct {
	models.Model
	CaseNo       string     `gorm:"size:32;uniqueIndex;not null" json:"caseNo"`
	Source       int        `gorm:"type:tinyint;not null" json:"source"`
	SourceID     uint64     `gorm:"default:0" json:"sourceId"`
	TargetType   int        `gorm:"type:tinyint;not null;index" json:"targetType"`
	TargetID     string     `gorm:"size:128;not null;index" json:"targetId"`
	Title        string     `gorm:"size:200;not null" json:"title"`
	Description  string     `gorm:"type:text" json:"description"`
	Priority     int        `gorm:"type:tinyint;default:1" json:"priority"`
	Status       int        `gorm:"type:tinyint;not null;default:1;index" json:"status"`
	HandlerID    string     `gorm:"size:64;index" json:"handlerId"`
	HandleRemark string     `gorm:"type:text" json:"handleRemark"`
	HandleTime   *time.Time `json:"handleTime"`
	ActionsTaken string     `gorm:"type:text" json:"actionsTaken"`
}
