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

package circle_models

import "beaver/common/models"

// CircleInviteModel 圈子分享邀请凭证
type CircleInviteModel struct {
	models.Model
	Token      string `gorm:"column:token;size:64;uniqueIndex;not null" json:"token"` // 对外邀请凭证
	CircleID   string `gorm:"column:circle_id;size:64;not null;index" json:"circleId"`
	CreatorID  string `gorm:"column:creator_id;size:64;not null;index" json:"creatorId"`
	ExpireAt   int64  `gorm:"column:expire_at;not null;index" json:"expireAt"` // Unix 秒
	MaxUses    int64  `gorm:"column:max_uses;not null;default:0" json:"maxUses"` // 0=不限
	UsedCount  int64  `gorm:"column:used_count;not null;default:0" json:"usedCount"`
	Status     int8   `gorm:"column:status;not null;default:1;index" json:"status"` // 1有效 2吊销 3用尽
}
