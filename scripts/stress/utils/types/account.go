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

/*
 * 共享数据类型：测试账号与好友关系，供 setup 脚本和压测脚本共用。
 *
 * 纯类型定义，独立成 package types，不依赖任何业务逻辑。
 */

package types

// Account 测试账号，序列化到 accounts.json 供所有压测脚本共用
type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   string `json:"userId,omitempty"`
	Token    string `json:"token,omitempty"`
}

// TestPassword 测试账号统一密码
const TestPassword = "Test@123456"

// Friendship 一对好友关系 + 自动产生的会话ID，供 chat 压测取 conversationId。
// 序列化到 friends.json，由 setup/friends 产出。
type Friendship struct {
	UserID         string `json:"userId"`
	FriendID       string `json:"friendId"`
	ConversationID string `json:"conversationId"`
}
