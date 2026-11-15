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

// beaver-moment S2S 出入参类型
//
// 手写 JSON DTO：字段名与 json tag 与对端 proto 的 JSON 编码逐字段一致，
// 不依赖 protobuf 运行时（无 protoimpl / protoreflect / state）。
package moment

type DeleteMomentCommentsReq struct {
	CommentIds []string `json:"comment_ids,omitempty"`
}

type DeleteMomentCommentsRes struct{}

type ListMomentCommentsReq struct {
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"page_size,omitempty"`
	MomentId string `json:"moment_id,omitempty"`
}

type ListMomentCommentsRes struct {
	Total int64                `json:"total,omitempty"`
	List  []*MomentCommentItem `json:"list,omitempty"`
}

type ListMomentsReq struct {
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"page_size,omitempty"`
	UserId   string `json:"user_id,omitempty"`
	Keywords string `json:"keywords,omitempty"`
	MomentId string `json:"moment_id,omitempty"`
}

type ListMomentsRes struct {
	Total int64         `json:"total,omitempty"`
	List  []*MomentItem `json:"list,omitempty"`
}

type MomentCommentItem struct {
	CommentId string `json:"comment_id,omitempty"`
	MomentId  string `json:"moment_id,omitempty"`
	UserId    string `json:"user_id,omitempty"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type MomentFileItem struct {
	FileKey string `json:"file_key,omitempty"`
}

type MomentItem struct {
	MomentId  string            `json:"moment_id,omitempty"`
	UserId    string            `json:"user_id,omitempty"`
	Content   string            `json:"content,omitempty"`
	Files     []*MomentFileItem `json:"files,omitempty"`
	IsDeleted bool              `json:"is_deleted,omitempty"`
	CreatedAt string            `json:"created_at,omitempty"`
	UpdatedAt string            `json:"updated_at,omitempty"`
}

type UpdateMomentReq struct {
	MomentId  string `json:"moment_id,omitempty"`
	IsDeleted *bool  `json:"is_deleted,omitempty"`
}

type UpdateMomentRes struct{}
