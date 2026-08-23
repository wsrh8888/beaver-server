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

package call_models

import (
	"beaver/common/models"
	"time"
)

/*
SessionStatus 业务逻辑说明：
1. 初始状态/进行中 (1): A发起呼叫时默认为1；只要有1人同意并留在通话中，状态始终为1。
2. 已结束 (3): 通话正常进行后，所有人均已退出房间。
3. 未接听 (4): 呼叫发出后，在规定时间内没有任何人回应（全部超时）。
4. 拒接 (5): 所有的受邀者都明确点击了“拒绝”。
*/

// SessionStatus 会话状态
type SessionStatus int8

const (
	SessionStatusCalling  SessionStatus = 1 // 待接听/进行中
	SessionStatusEnded    SessionStatus = 3 // 已结束
	SessionStatusMissed   SessionStatus = 4 // 未接听 (全部超时)
	SessionStatusRejected SessionStatus = 5 // 拒接 (全部拒绝)
)

// CallSession 通话主表（对应一次呼叫流程）
type CallSession struct {
	models.Model
	RoomID         string        `gorm:"type:varchar(64);uniqueIndex;not null;comment:LiveKit房间ID/业务唯一ID" json:"room_id"`
	CallerID       string        `gorm:"type:varchar(64);index;not null;comment:发起者ID" json:"caller_id"`
	ConversationID string        `gorm:"type:varchar(64);index;comment:会话ID" json:"conversation_id"`
	CallType       int8          `gorm:"type:tinyint;default:1;comment:通话类型:1-私聊,2-群聊" json:"call_type"`
	Status         SessionStatus `gorm:"type:tinyint;default:1;index;comment:状态:1-待接听/进行中,3-已结束,4-未接听,5-拒接" json:"status"`
	StartTime      *time.Time    `gorm:"type:datetime;comment:接通时间" json:"start_time"`
	EndTime        *time.Time    `gorm:"type:datetime;comment:挂断时间" json:"end_time"`
	Duration       int32         `gorm:"type:int;default:0;comment:通话时长(秒)" json:"duration"`
}
