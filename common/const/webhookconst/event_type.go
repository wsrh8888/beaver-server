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

package webhookconst

// Webhook 事件类型常量
const (
	// Bot 相关事件
	EventBotMessageReceive = "bot.message.receive" // Bot 收到消息

	// 消息相关事件
	EventAfterSendMsg  = "after.send_msg"  // 消息发送后
	EventBeforeSendMsg = "before.send_msg" // 消息发送前(可拦截)
	EventMessageRecall = "message.recall"  // 消息撤回

	// 好友相关事件
	EventAfterAddFriend    = "after.add_friend"    // 加好友后
	EventAfterDeleteFriend = "after.delete_friend" // 删除好友后

	// 群组相关事件
	EventAfterJoinGroup  = "after.join_group"  // 加群后
	EventAfterLeaveGroup = "after.leave_group" // 退群后
	EventGroupDissolve   = "group.dissolve"    // 群组解散

	// 用户相关事件
	EventUserRegister = "user.register" // 用户注册
	EventUserLogin    = "user.login"    // 用户登录
)
