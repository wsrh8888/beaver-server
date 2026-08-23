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

import "beaver/common/models"

// AdminOperationLog 管理员操作审计日志
type AdminOperationLog struct {
	models.Model
	OperatorID   string `gorm:"size:64;index;not null" json:"operatorId"`
	Action       string `gorm:"size:64;index;not null" json:"action"`
	TargetType   string `gorm:"size:32;index" json:"targetType"`
	TargetID     string `gorm:"size:128;index" json:"targetId"`
	CaseID       uint64 `gorm:"default:0;index" json:"caseId"`
	Detail       string `gorm:"type:text" json:"detail"`
	Result       string `gorm:"size:32;not null" json:"result"`
	ErrorMessage string `gorm:"type:text" json:"errorMessage"`
}
