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

package datasync_models

import (
	"beaver/common/models"
)

// DatasyncModel 数据同步模型
type DatasyncModel struct {
	models.Model
	DataType     string `gorm:"size:32;unique" json:"dataType"`              // 数据类型：users/friends/groups/chats/conversations（唯一键）
	LastSeq      int64  `gorm:"default:0" json:"lastSeq"`                    // 最后同步的序列号（消息用）或版本号（基础数据用）
	LastSyncTime int64  `gorm:"default:0" json:"lastSyncTime"`               // 最后同步的时间戳
	SyncStatus   string `gorm:"size:16;default:'pending'" json:"syncStatus"` // 同步状态：pending/syncing/completed/failed
}
