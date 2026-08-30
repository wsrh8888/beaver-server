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

package openevent

import "fmt"

// 智能机器人 Robot 平台 Webhook 事件类型
const (
	EventIMMessageReceive       = "im.message.receive"
	EventIMMessageReceiveGroup  = "im.message.receive.group_at"
	EventIMBotFollowed          = "im.bot.followed"
	EventIMBotUnfollowed        = "im.bot.unfollowed"
	EventIMChatMemberBotAdded   = "im.chat.member.bot.added"
	EventIMChatMemberBotRemoved = "im.chat.member.bot.removed"
)

func ValidateRobotEventType(eventType string) error {
	switch eventType {
	case EventIMMessageReceive, EventIMMessageReceiveGroup,
		EventIMBotFollowed, EventIMBotUnfollowed,
		EventIMChatMemberBotAdded, EventIMChatMemberBotRemoved:
		return nil
	default:
		return fmt.Errorf("不支持的事件类型: %s", eventType)
	}
}
