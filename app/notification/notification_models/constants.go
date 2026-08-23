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

package notification_models

// 事件类型常量，便于统一引用和维护
const (
	// 动态相关
	EventTypeMomentLike         = "moment_like"
	EventTypeMomentUnlike       = "moment_unlike"
	EventTypeMomentComment      = "moment_comment"
	EventTypeMomentCommentReply = "moment_comment_reply"

	// 好友相关
	EventTypeFriendRequest       = "friend_request"
	EventTypeFriendRequestAccept = "friend_request_accept"
	EventTypeFriendRequestReject = "friend_request_reject"

	// 群相关
	EventTypeGroupJoinRequest = "group_join_request"
	EventTypeGroupLeft        = "group_left"

)

// 分类常量
const (
	CategoryMoment = "moment" // 动态
	CategorySocial = "social" // 社交
	CategoryGroup  = "group"  // 群组
)

// 目标类型常量（事件指向的实体类型）
const (
	TargetTypeMoment        = "moment"
	TargetTypeMomentComment = "moment_comment"
	TargetTypeGroup         = "group"
	TargetTypeUser          = "user"
)

// 版本号作用域（用于 VersionGen 约定）
const (
	// NotificationEvent：全局递增版本
	VersionScopeEventGlobal = "notification_events"
	// NotificationInbox：按用户递增版本（per user）
	VersionScopeInboxPerUser = "notification_inboxes"
	// NotificationRead：按用户+分类递增版本（per user per category）
	VersionScopeCursorPerUser = "notification_reads"
)
