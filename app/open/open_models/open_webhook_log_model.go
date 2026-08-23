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

import "gorm.io/gorm"

// OpenWebhookLog Webhook 推送日志
type OpenWebhookLog struct {
	gorm.Model
	SubscriptionID uint   `gorm:"index;comment:订阅ID"`
	ConfigID       string `gorm:"type:varchar(64);index;comment:配置ID(兼容)"`
	AppID          string `gorm:"type:varchar(64);index;comment:应用ID"`
	EventID        string `gorm:"type:varchar(64);index;comment:事件ID"`
	EventType      string `gorm:"type:varchar(100);comment:事件类型"`
	TargetURL      string `gorm:"type:varchar(512);comment:目标URL"`
	HTTPStatus     int    `gorm:"comment:HTTP状态码"`
	LatencyMs      int64  `gorm:"comment:耗时毫秒"`
	RetryCount     int    `gorm:"comment:重试次数"`
	ErrorMessage   string `gorm:"type:varchar(512);comment:错误信息"`
	Status         int    `gorm:"type:tinyint;default:0;comment:1成功 0失败"`
}

func (OpenWebhookLog) TableName() string {
	return "open_webhook_logs"
}
