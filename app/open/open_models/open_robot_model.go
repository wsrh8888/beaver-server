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
	"gorm.io/gorm"
)

// OpenRobot 应用智能机器人（一个应用对应一个机器人 IM 用户）
//
// 职责边界：本表**只表达归属关系** —— 哪个应用对应哪个 IM 用户、是否启用。
//
// 机器人的一切「对外行为」都不在这里：
//   昵称 / 头像        → UserModel（RobotID 即 UserModel.UserID）
//   欢迎语 / 命令前缀   → 机器人自己（agent 侧）
//   是否在单聊/群聊响应  → 机器人自己判断（平台只负责把消息送达，不替它决定该不该回）
type OpenRobot struct {
	gorm.Model
	AppID   string `gorm:"type:varchar(64);uniqueIndex;not null;comment:应用ID"`
	RobotID string `gorm:"column:robot_user_id;type:varchar(64);uniqueIndex;comment:Robot IM 用户ID"`
	Status  int    `gorm:"type:tinyint;default:1;comment:1启用 0禁用"`
}
