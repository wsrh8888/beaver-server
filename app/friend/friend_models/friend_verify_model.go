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

package friend_models

import (
	"beaver/common/models"
)

/**
 * @description: 好友验证
 */
type FriendVerifyModel struct {
	models.Model
	VerifyID   string `gorm:"column:verify_id;size:64;uniqueIndex" json:"verifyId"`
	SendUserID string `gorm:"size:64;index" json:"sendUserId"` // 使用 VARCHAR(64)
	RevUserID  string `gorm:"size:64;index" json:"revUserId"`  // 使用 VARCHAR(64)
	SendStatus int8   `json:"sendStatus"`                      // 发起方状态 0:未处理 1:已通过 2:已拒绝 3: 忽略 4:删除
	RevStatus  int8   `json:"revStatus"`                       // 接收方状态 0:未处理 1:已通过 2:已拒绝 3: 忽略 4:删除
	Message    string `gorm:"size: 128" json:"message"`        // 附加消息
	Source     string `gorm:"size: 32" json:"source"`          // 添加好友来源：qrcode/search/group/recommend
	Version    int64  `gorm:"not null;default:0;index"`        // 序列号（用于数据同步）
}
