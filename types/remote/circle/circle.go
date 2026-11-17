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

// beaver-circle S2S 出入参类型
//
// 手写 JSON DTO：字段名与 json tag 与对端 proto 的 JSON 编码逐字段一致，
// 不依赖 protobuf 运行时（无 protoimpl / protoreflect / state）。
package circle

// DeleteCommentApi 删除评论：Req + Res。
type (
	DeleteCommentReq struct {
		CommentId string `json:"comment_id,omitempty"`
	}

	DeleteCommentRes struct{}
)

// DeletePostApi 删除帖子：Req + Res。
type (
	DeletePostReq struct {
		PostId string `json:"post_id,omitempty"`
	}

	DeletePostRes struct{}
)

// GetCircleListApi 列圈子：Req + Res + 本接口专属子类型。
type (
	GetCircleListReq struct {
		Page     int32  `json:"page,omitempty"`
		PageSize int32  `json:"page_size,omitempty"`
		UserId   string `json:"user_id,omitempty"`
		Keywords string `json:"keywords,omitempty"`
		CircleId string `json:"circle_id,omitempty"`
	}

	GetCircleListRes struct {
		Total int64                      `json:"total,omitempty"`
		List  []*GetCircleListCircleItem `json:"list,omitempty"`
	}

	// GetCircleListCircleItem 圈子项。
	GetCircleListCircleItem struct {
		CircleId    string `json:"circle_id,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Avatar      string `json:"avatar,omitempty"`
		CreatorId   string `json:"creator_id,omitempty"`
		JoinType    int32  `json:"join_type,omitempty"`
		MemberCount int64  `json:"member_count,omitempty"`
		PostCount   int64  `json:"post_count,omitempty"`
		IsDeleted   bool   `json:"is_deleted,omitempty"`
		CreatedAt   string `json:"created_at,omitempty"`
		UpdatedAt   string `json:"updated_at,omitempty"`
	}
)

// GetCircleMembersApi 列圈子成员：Req + Res + 本接口专属子类型。
type (
	GetCircleMembersReq struct {
		CircleId string `json:"circle_id,omitempty"`
		Page     int32  `json:"page,omitempty"`
		PageSize int32  `json:"page_size,omitempty"`
	}

	GetCircleMembersRes struct {
		Total int64                         `json:"total,omitempty"`
		List  []*GetCircleMembersMemberItem `json:"list,omitempty"`
	}

	// GetCircleMembersMemberItem 圈子成员项。
	GetCircleMembersMemberItem struct {
		CircleId string `json:"circle_id,omitempty"`
		UserId   string `json:"user_id,omitempty"`
		Role     int32  `json:"role,omitempty"`
	}
)

// GetCircleVersionsApi 列圈子版本变更：Req + Res + 本接口专属子类型。
type (
	GetCircleVersionsReq struct {
		UserId  string `json:"user_id,omitempty"`
		Version int64  `json:"version,omitempty"`
	}

	GetCircleVersionsRes struct {
		List []*GetCircleVersionsVersionItem `json:"list,omitempty"`
	}

	// GetCircleVersionsVersionItem 圈子版本项。
	GetCircleVersionsVersionItem struct {
		CircleId    string `json:"circle_id,omitempty"`
		Name        string `json:"name,omitempty"`
		Avatar      string `json:"avatar,omitempty"`
		MemberCount int64  `json:"member_count,omitempty"`
		Role        int32  `json:"role,omitempty"`
		Version     int64  `json:"version,omitempty"`
	}
)

// GetUserCircleIDsApi 取用户加入的圈子 ID：Req + Res。
type (
	GetUserCircleIDsReq struct {
		UserId string `json:"user_id,omitempty"`
	}

	GetUserCircleIDsRes struct {
		CircleIds []string `json:"circle_ids,omitempty"`
	}
)

// GetUserCircleRoleApi 取用户圈子角色：Req + Res。
type (
	GetUserCircleRoleReq struct {
		CircleId string `json:"circle_id,omitempty"`
		UserId   string `json:"user_id,omitempty"`
	}

	GetUserCircleRoleRes struct {
		Role int32 `json:"role,omitempty"`
	}
)

// ListCommentsApi 列帖子评论：Req + Res + 本接口专属子类型。
type (
	ListCommentsReq struct {
		PostId   string `json:"post_id,omitempty"`
		Page     int32  `json:"page,omitempty"`
		PageSize int32  `json:"page_size,omitempty"`
	}

	ListCommentsRes struct {
		Total int64                      `json:"total,omitempty"`
		List  []*ListCommentsCommentItem `json:"list,omitempty"`
	}

	// ListCommentsCommentItem 评论项。
	ListCommentsCommentItem struct {
		CommentId string `json:"comment_id,omitempty"`
		PostId    string `json:"post_id,omitempty"`
		UserId    string `json:"user_id,omitempty"`
		Content   string `json:"content,omitempty"`
		IsDeleted bool   `json:"is_deleted,omitempty"`
		CreatedAt string `json:"created_at,omitempty"`
	}
)

// ListPostsApi 列圈子帖子：Req + Res + 本接口专属子类型。
type (
	ListPostsReq struct {
		CircleId string `json:"circle_id,omitempty"`
		Page     int32  `json:"page,omitempty"`
		PageSize int32  `json:"page_size,omitempty"`
	}

	ListPostsRes struct {
		Total int64                `json:"total,omitempty"`
		List  []*ListPostsPostItem `json:"list,omitempty"`
	}

	// ListPostsPostItem 帖子项。
	ListPostsPostItem struct {
		PostId       string `json:"post_id,omitempty"`
		CircleId     string `json:"circle_id,omitempty"`
		UserId       string `json:"user_id,omitempty"`
		Content      string `json:"content,omitempty"`
		IsDeleted    bool   `json:"is_deleted,omitempty"`
		CreatedAt    string `json:"created_at,omitempty"`
		CommentCount int64  `json:"comment_count,omitempty"`
		LikeCount    int64  `json:"like_count,omitempty"`
	}
)

// RemoveCircleMembersApi 移除圈子成员：Req + Res。
type (
	RemoveCircleMembersReq struct {
		CircleId string   `json:"circle_id,omitempty"`
		UserIds  []string `json:"user_ids,omitempty"`
	}

	RemoveCircleMembersRes struct{}
)

// UpdateCircleApi 更新圈子：Req + Res。
type (
	UpdateCircleReq struct {
		CircleId  string `json:"circle_id,omitempty"`
		IsDeleted *bool  `json:"is_deleted,omitempty"`
	}

	UpdateCircleRes struct{}
)
