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

// OpenAppEventSubscription 应用事件订阅（一种 eventType 一条，对标飞书）
type OpenAppEventSubscription struct {
	gorm.Model
	AppID       string `gorm:"type:varchar(64);uniqueIndex:uk_app_event;not null;comment:应用ID"`
	EventType   string `gorm:"type:varchar(100);uniqueIndex:uk_app_event;not null;comment:事件类型"`
	CallbackURL string `gorm:"type:varchar(512);not null;comment:Webhook URL"`
	Secret      string `gorm:"type:varchar(128);comment:签名密钥"`
	Status      int    `gorm:"type:tinyint;default:1;comment:1启用 0禁用"`
	VerifyStatus int   `gorm:"type:tinyint;default:0;comment:0待验证 1已通过 2失败"`
	LastVerifiedAt *time.Time `gorm:"comment:上次验证时间"`
	LastError   string `gorm:"type:varchar(512);comment:上次失败原因"`
	RetryCount  int    `gorm:"type:int;default:3;comment:重试次数"`
	Timeout     int    `gorm:"type:int;default:5;comment:超时秒"`
}
