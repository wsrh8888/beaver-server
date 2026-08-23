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

package constants

// ScopeType 权限范围类型
type ScopeType string

const (
	// 基础权限
	ScopeUserProfileRead ScopeType = "user_profile:read" // 读取用户基本信息
	ScopeUserAvatarRead  ScopeType = "user_avatar:read"  // 读取用户头像
	ScopeUserEmailRead   ScopeType = "user_email:read"   // 读取用户邮箱

	// 好友权限
	ScopeFriendListRead ScopeType = "friend:list:read" // 读取好友列表
	ScopeFriendInfoRead ScopeType = "friend:info:read" // 读取好友信息

	// 群组权限
	ScopeGroupListRead   ScopeType = "group:list:read"   // 读取群组列表
	ScopeGroupInfoRead   ScopeType = "group:info:read"   // 读取群组信息
	ScopeGroupMemberRead ScopeType = "group:member:read" // 读取群组成员

	// 消息权限
	ScopeMessageSend ScopeType = "message:send" // 发送消息
	ScopeMessageRead ScopeType = "message:read" // 读取消息

	// 文件权限
	ScopeFileUpload   ScopeType = "file:upload"   // 上传文件
	ScopeFileDownload ScopeType = "file:download" // 下载文件
)

// AllScopes 所有可用的权限列表
var AllScopes = []ScopeType{
	ScopeUserProfileRead,
	ScopeUserAvatarRead,
	ScopeUserEmailRead,
	ScopeFriendListRead,
	ScopeFriendInfoRead,
	ScopeGroupListRead,
	ScopeGroupInfoRead,
	ScopeGroupMemberRead,
	ScopeMessageSend,
	ScopeMessageRead,
	ScopeFileUpload,
	ScopeFileDownload,
}

// DefaultScopes 默认权限（创建应用时自动授予）
var DefaultScopes = []ScopeType{
	ScopeUserProfileRead,
	ScopeUserAvatarRead,
}

// ScopeDescription 权限描述
var ScopeDescription = map[ScopeType]string{
	ScopeUserProfileRead: "读取用户基本信息（昵称、头像等）",
	ScopeUserAvatarRead:  "读取用户头像",
	ScopeUserEmailRead:   "读取用户邮箱地址",
	ScopeFriendListRead:  "读取用户的好友列表",
	ScopeFriendInfoRead:  "读取好友详细信息",
	ScopeGroupListRead:   "读取用户加入的群组列表",
	ScopeGroupInfoRead:   "读取群组基本信息",
	ScopeGroupMemberRead: "读取群组成员列表",
	ScopeMessageSend:     "以应用身份发送消息",
	ScopeMessageRead:     "读取应用相关的消息",
	ScopeFileUpload:      "上传文件到服务器",
	ScopeFileDownload:    "下载文件",
}
