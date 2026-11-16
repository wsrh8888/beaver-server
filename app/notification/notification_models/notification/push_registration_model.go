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

package notification_models

import (
	"beaver/common/models"
)

// PushRegistrationModel 离线推送 Token 注册表（notification 库）
type PushRegistrationModel struct {
	models.Model
	UserID       string `gorm:"size:64;not null;uniqueIndex:idx_push_user_device" json:"userId"`
	DeviceID     string `gorm:"size:128;not null;uniqueIndex:idx_push_user_device" json:"deviceId"`
	PushToken    string `gorm:"size:512;not null" json:"-"`
	PushPlatform string `gorm:"size:16;not null;index" json:"pushPlatform"` // fcm | apns
	Enabled      bool   `gorm:"not null;default:true;index" json:"enabled"`
}
