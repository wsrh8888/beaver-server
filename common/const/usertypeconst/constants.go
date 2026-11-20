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

// Package usertypeconst 存放 UserType / Source 公共枚举，避免各域 import account_models。
package usertypeconst

const (
	UserTypeNormal  int8 = 1 // 普通用户
	UserTypeBot     int8 = 2 // Bot 机器人（外接/推送机器人），工号前缀 B
	UserTypeDigital int8 = 3 // Agent 数字伙伴（智能体），工号前缀 A
)

// 用户注册来源（account 域 user.source 的取值）
const (
	SourcePhone int32 = 1 // 手机号注册
	SourceEmail int32 = 2 // 邮箱注册
	SourceGroup int32 = 3 // 群内创建（机器人）
)
