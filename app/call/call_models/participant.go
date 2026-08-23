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
ParticipantStatus 业务逻辑说明：
1. 待接听 (1): 正在被呼叫，手机正在振铃中。
2. 已接听 (2): 用户点击接听，成功进入音视频房间。
3. 拒绝 (3): 振铃阶段用户手动点击“挂断”或“拒绝”。
4. 超时未接 (4): 振铃超时，用户未操作，系统自动触发取消。
5. 已退出 (5): 曾经进场，后续正常离开、手动挂断或掉线。
*/

// ParticipantStatus 参与者状态
type ParticipantStatus int8

const (
	ParticipantStatusCalling  ParticipantStatus = 1 // 待接听
	ParticipantStatusJoined   ParticipantStatus = 2 // 已接听
	ParticipantStatusRejected ParticipantStatus = 3 // 拒绝
	ParticipantStatusTimeout  ParticipantStatus = 4 // 超时未接
	ParticipantStatusLeft     ParticipantStatus = 5 // 已退出
)

// CallParticipant 通话参与者表
type CallParticipant struct {
	models.Model
	RoomID string `gorm:"type:varchar(64);index:idx_room_user;not null;comment:关联RoomID" json:"room_id"`
	UserID string `gorm:"type:varchar(64);index:idx_room_user;not null;comment:用户ID" json:"user_id"`
	// 核心行为状态
	Status ParticipantStatus `gorm:"type:tinyint;default:1;comment:状态:1-待接听,2-已接听,3-拒绝,4-超时,5-挂断" json:"status"`
	Role   int8              `gorm:"type:tinyint;default:1;comment:角色:1-发起者,2-受邀者" json:"role"`

	JoinTime  *time.Time `gorm:"type:datetime;comment:加入时间" json:"join_time"`
	LeaveTime *time.Time `gorm:"type:datetime;comment:离开时间" json:"leave_time"`

	// 扩展信息（用于质量监控，大厂通常会有）
	DeviceInfo string `gorm:"type:varchar(255);comment:设备信息(iOS/Android)" json:"device_info"`
}
