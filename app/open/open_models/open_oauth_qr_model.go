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

// OpenOAuthQrCode 扫码登录记录表
type OpenOAuthQrCode struct {
	ID        uint           `gorm:"primarykey"`
	SceneID   string         `gorm:"column:scene_id;type:varchar(64);not null;index"`
	AppID     string         `gorm:"column:app_id;type:varchar(64);not null;index"`
	UserID    string         `gorm:"column:user_id;type:varchar(64);default:''"`
	Status    int            `gorm:"column:status;type:tinyint;not null;default:0;comment:0-等待扫码,1-已扫码,2-已确认,3-已取消,4-已过期"`
	ExpiresAt time.Time      `gorm:"column:expires_at;type:datetime;not null"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;index"`
}

func (OpenOAuthQrCode) TableName() string {
	return "open_oauth_qr_codes"
}
