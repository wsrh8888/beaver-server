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

package open_models

import (
	"time"

	"gorm.io/gorm"
)

// OpenOAuthQrCode 扫码登录记录（一次扫码会话一条，用完可清理）
type OpenOAuthQrCode struct {
	gorm.Model
	SceneID   string    `gorm:"column:scene_id;type:varchar(64);not null;index"`
	AppID     string    `gorm:"column:app_id;type:varchar(64);not null;index"`
	UserID    string    `gorm:"column:user_id;type:varchar(64);comment:扫码用户，空表示尚未扫码"`
	Status    int       `gorm:"column:status;type:tinyint;not null;default:0;comment:0等待扫码 1已扫码 2已确认 3已取消 4已过期"`
	ExpiresAt time.Time `gorm:"column:expires_at;type:datetime;not null"`
}

func (OpenOAuthQrCode) TableName() string {
	return "open_oauth_qr_codes"
}
