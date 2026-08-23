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

package platform_models

import "beaver/common/models"

// UpdateReleasePolicy 架构发版策略（正式版 + 比例灰度 + 最低强更）
type UpdateReleasePolicy struct {
	models.Model
	AppID            string `json:"appId" gorm:"type:varchar(64);index"`
	ArchitectureID   uint   `json:"architectureId" gorm:"uniqueIndex"`
	StableVersionID  uint   `json:"stableVersionId"`
	GrayVersionID    uint   `json:"grayVersionId"`
	RolloutPercent   uint   `json:"rolloutPercent"` // 0-100，命中灰度桶的用户使用 GrayVersionID
	MinVersion       string `json:"minVersion" gorm:"type:varchar(32)"`
	ForceUpdate      bool   `json:"forceUpdate"`
	IsActive         bool   `json:"isActive"`
}
