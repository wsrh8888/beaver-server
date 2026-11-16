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

// Package notifyconst 存放通知事件的分类/类型等公共常量。
// 这些常量被多域共享（notification 投递方、relation 事件生产方等），
// 下沉到 common 避免各域互相 import 对方的 models 包。
package notifyconst

// 事件类型
const (
	EventTypeMomentLike         = "moment_like"
	EventTypeMomentUnlike       = "moment_unlike"
	EventTypeMomentComment      = "moment_comment"
	EventTypeMomentCommentReply = "moment_comment_reply"

	EventTypeFriendRequest       = "friend_request"
	EventTypeFriendRequestAccept = "friend_request_accept"
	EventTypeFriendRequestReject = "friend_request_reject"

	EventTypeGroupJoinRequest = "group_join_request"
	EventTypeGroupLeft        = "group_left"
)

// 分类
const (
	CategoryMoment = "moment"
	CategorySocial = "social"
	CategoryGroup  = "group"
)

// 目标类型
const (
	TargetTypeMoment        = "moment"
	TargetTypeMomentComment = "moment_comment"
	TargetTypeGroup         = "group"
	TargetTypeUser          = "user"
)

// 版本号作用域
const (
	VersionScopeEventGlobal   = "notification_events"
	VersionScopeInboxPerUser  = "notification_inboxes"
	VersionScopeCursorPerUser = "notification_reads"
)
