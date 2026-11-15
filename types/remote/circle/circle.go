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

type CircleItem struct {
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

type CircleMemberItem struct {
	CircleId string `json:"circle_id,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	Role     int32  `json:"role,omitempty"`
}

type CircleVersionItem struct {
	CircleId    string `json:"circle_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	MemberCount int64  `json:"member_count,omitempty"`
	Role        int32  `json:"role,omitempty"`
	Version     int64  `json:"version,omitempty"`
}

type CommentItem struct {
	CommentId string `json:"comment_id,omitempty"`
	PostId    string `json:"post_id,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	Content   string `json:"content,omitempty"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

type DeleteCommentReq struct {
	CommentId string `json:"comment_id,omitempty"`
}

type DeleteCommentRes struct{}

type DeletePostReq struct {
	PostId string `json:"post_id,omitempty"`
}

type DeletePostRes struct{}

type GetCircleListReq struct {
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"page_size,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	Keywords string `json:"keywords,omitempty"`
	CircleId string `json:"circle_id,omitempty"`
}

type GetCircleListRes struct {
	Total int64         `json:"total,omitempty"`
	List  []*CircleItem `json:"list,omitempty"`
}

type GetCircleMembersReq struct {
	CircleId string `json:"circle_id,omitempty"`
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"page_size,omitempty"`
}

type GetCircleMembersRes struct {
	Total int64               `json:"total,omitempty"`
	List  []*CircleMemberItem `json:"list,omitempty"`
}

type GetCircleVersionsReq struct {
	UserId  string `json:"user_id,omitempty"`
	Version int64  `json:"version,omitempty"`
}

type GetCircleVersionsRes struct {
	List []*CircleVersionItem `json:"list,omitempty"`
}

type GetUserCircleIDsReq struct {
	UserId string `json:"user_id,omitempty"`
}

type GetUserCircleIDsRes struct {
	CircleIds []string `json:"circle_ids,omitempty"`
}

type GetUserCircleRoleReq struct {
	CircleId string `json:"circle_id,omitempty"`
	UserId   string `json:"user_id,omitempty"`
}

type GetUserCircleRoleRes struct {
	Role int32 `json:"role,omitempty"`
}

type ListCommentsReq struct {
	PostId   string `json:"post_id,omitempty"`
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"page_size,omitempty"`
}

type ListCommentsRes struct {
	Total int64          `json:"total,omitempty"`
	List  []*CommentItem `json:"list,omitempty"`
}

type ListPostsReq struct {
	CircleId string `json:"circle_id,omitempty"`
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"page_size,omitempty"`
}

type ListPostsRes struct {
	Total int64       `json:"total,omitempty"`
	List  []*PostItem `json:"list,omitempty"`
}

type PostItem struct {
	PostId       string `json:"post_id,omitempty"`
	CircleId     string `json:"circle_id,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	Content      string `json:"content,omitempty"`
	IsDeleted    bool   `json:"is_deleted,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	CommentCount int64  `json:"comment_count,omitempty"`
	LikeCount    int64  `json:"like_count,omitempty"`
}

type RemoveCircleMembersReq struct {
	CircleId string   `json:"circle_id,omitempty"`
	UserIds  []string `json:"user_ids,omitempty"`
}

type RemoveCircleMembersRes struct{}

type UpdateCircleReq struct {
	CircleId  string `json:"circle_id,omitempty"`
	IsDeleted *bool  `json:"is_deleted,omitempty"`
}

type UpdateCircleRes struct{}
