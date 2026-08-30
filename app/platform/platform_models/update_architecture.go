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

import (
	"beaver/common/models"
)

// UpdateArchitecture 架构信息表
type UpdateArchitecture struct {
	models.Model
	AppID       string `json:"appId" gorm:"type:varchar(64);index"`
	PlatformID  uint   `json:"platformId"`
	ArchID      uint   `json:"archId"`
	Description string `json:"description"`
	IsActive    bool   `json:"isActive"`
}

const (
	PlatformWindows   uint = 1
	PlatformMacOS     uint = 2
	PlatformIOS       uint = 3
	PlatformAndroid   uint = 4
	PlatformHarmonyOS uint = 5
)

const (
	H5        uint = 0
	WinX64    uint = 1
	WinArm64  uint = 2
	MacIntel  uint = 3
	MacApple  uint = 4
	IOS       uint = 5
	Android   uint = 6
	HarmonyOS uint = 7
)
