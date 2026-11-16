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

import "beaver/common/const/notifyconst"

// 通知事件常量统一定义在 common/const/notifyconst，这里做别名转发，
// 保持 notification 域内 notification_models.Xxx 的既有引用不变。
const (
	EventTypeMomentLike         = notifyconst.EventTypeMomentLike
	EventTypeMomentUnlike       = notifyconst.EventTypeMomentUnlike
	EventTypeMomentComment      = notifyconst.EventTypeMomentComment
	EventTypeMomentCommentReply = notifyconst.EventTypeMomentCommentReply

	EventTypeFriendRequest       = notifyconst.EventTypeFriendRequest
	EventTypeFriendRequestAccept = notifyconst.EventTypeFriendRequestAccept
	EventTypeFriendRequestReject = notifyconst.EventTypeFriendRequestReject

	EventTypeGroupJoinRequest = notifyconst.EventTypeGroupJoinRequest
	EventTypeGroupLeft        = notifyconst.EventTypeGroupLeft
)

const (
	CategoryMoment = notifyconst.CategoryMoment
	CategorySocial = notifyconst.CategorySocial
	CategoryGroup  = notifyconst.CategoryGroup
)

const (
	TargetTypeMoment        = notifyconst.TargetTypeMoment
	TargetTypeMomentComment = notifyconst.TargetTypeMomentComment
	TargetTypeGroup         = notifyconst.TargetTypeGroup
	TargetTypeUser          = notifyconst.TargetTypeUser
)

const (
	VersionScopeEventGlobal   = notifyconst.VersionScopeEventGlobal
	VersionScopeInboxPerUser  = notifyconst.VersionScopeInboxPerUser
	VersionScopeCursorPerUser = notifyconst.VersionScopeCursorPerUser
)
