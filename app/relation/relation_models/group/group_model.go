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

package group_models

import (
	"beaver/common/models"
	"time"
)

// 群组独立递增版本
type GroupModel struct {
	models.Model
	GroupID   string     `gorm:"size:64;unique;index" json:"groupId"`
	Type      int8       `gorm:"default:1" json:"type"`                                                                                                   // 群类型：1正常群 2讨论组 ...
	Title     string     `gorm:"size:32;index" json:"title"`                                                                                              // 群名
	Avatar    string     `gorm:"size:256;default:https://server.wsrh8888.com/beaver/api/file/preview/a8ba5d19ea54a91aec17dec0ad5000e6.png" json:"avatar"` // 群头像文件名
	CreatorID string     `gorm:"size:64;index" json:"creatorId"`                                                                                          // 创建者ID
	Notice    string     `gorm:"type:text" json:"notice"`                                                                                                 // 当前公告内容
	JoinType  int8       `gorm:"not null;default:0" json:"joinType"`                                                                                      // 0自由加入 1需审批 2不可加入
	IsMuteAll bool       `gorm:"not null;default:false" json:"isMuteAll"`                                                                                 // 是否全员禁言
	MuteAllAt *time.Time `json:"muteAllAt"`                                                                                                               // 开启全员禁言的时间
	Status    int8       `gorm:"default:1" json:"status"`                                                                                                 // 群状态：1正常 2冻结 3解散
	Version   int64      `gorm:"not null;default:0;index" json:"version"`                                                                                 // 群组版本号
}
