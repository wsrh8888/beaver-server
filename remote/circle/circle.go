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

// Package circle A 主体（开源核心）→ 收费模块圈子服务（beaver-server-social/app/social）的 S2S 接口代理。
//
// 跨仓不允许 etcd gRPC 直连（团队铁律），统一走对方网关 /internal/Circle.*（client_credentials 换票 + Bearer）。
// 结构对齐 remote/{open,ai,doc,moment}：enter.go 灌域名与凭据、remote.go 持有换票与发包 plumbing、
// 业务方法在本文件。Req/Res 放到 types/remote/circle。
package circle

import (
	"context"
	"net/http"

	"beaver/commonconfig"
	circletype "beaver/types/remote/circle"
)

// GetCircleList 获取圈子列表（支持按用户、关键词、ID 过滤）。
func GetCircleList(ctx context.Context, header http.Header, body *circletype.GetCircleListReq) (*circletype.GetCircleListRes, error) {
	return Post[circletype.GetCircleListRes](ctx, commonconfig.SocialUrl+"/internal/Circle.GetCircleList", header, body)
}

// UpdateCircle 更新圈子状态（解散等）。
func UpdateCircle(ctx context.Context, header http.Header, body *circletype.UpdateCircleReq) (*circletype.UpdateCircleRes, error) {
	return Post[circletype.UpdateCircleRes](ctx, commonconfig.SocialUrl+"/internal/Circle.UpdateCircle", header, body)
}

// GetCircleMembers 获取圈子成员列表。
func GetCircleMembers(ctx context.Context, header http.Header, body *circletype.GetCircleMembersReq) (*circletype.GetCircleMembersRes, error) {
	return Post[circletype.GetCircleMembersRes](ctx, commonconfig.SocialUrl+"/internal/Circle.GetCircleMembers", header, body)
}

// GetUserCircleRole 获取用户在圈子中的角色。
func GetUserCircleRole(ctx context.Context, header http.Header, body *circletype.GetUserCircleRoleReq) (*circletype.GetUserCircleRoleRes, error) {
	return Post[circletype.GetUserCircleRoleRes](ctx, commonconfig.SocialUrl+"/internal/Circle.GetUserCircleRole", header, body)
}

// GetUserCircleIDs 获取用户加入的所有圈子 ID。
func GetUserCircleIDs(ctx context.Context, header http.Header, body *circletype.GetUserCircleIDsReq) (*circletype.GetUserCircleIDsRes, error) {
	return Post[circletype.GetUserCircleIDsRes](ctx, commonconfig.SocialUrl+"/internal/Circle.GetUserCircleIDs", header, body)
}

// GetCircleVersions 获取圈子版本变更（数据同步）。
func GetCircleVersions(ctx context.Context, header http.Header, body *circletype.GetCircleVersionsReq) (*circletype.GetCircleVersionsRes, error) {
	return Post[circletype.GetCircleVersionsRes](ctx, commonconfig.SocialUrl+"/internal/Circle.GetCircleVersions", header, body)
}

// ListPosts 获取圈子帖子列表（管理后台）。
func ListPosts(ctx context.Context, header http.Header, body *circletype.ListPostsReq) (*circletype.ListPostsRes, error) {
	return Post[circletype.ListPostsRes](ctx, commonconfig.SocialUrl+"/internal/Circle.ListPosts", header, body)
}

// DeletePost 删除帖子（管理后台）。
func DeletePost(ctx context.Context, header http.Header, body *circletype.DeletePostReq) (*circletype.DeletePostRes, error) {
	return Post[circletype.DeletePostRes](ctx, commonconfig.SocialUrl+"/internal/Circle.DeletePost", header, body)
}

// ListComments 获取帖子评论列表（管理后台）。
func ListComments(ctx context.Context, header http.Header, body *circletype.ListCommentsReq) (*circletype.ListCommentsRes, error) {
	return Post[circletype.ListCommentsRes](ctx, commonconfig.SocialUrl+"/internal/Circle.ListComments", header, body)
}

// DeleteComment 删除评论（管理后台）。
func DeleteComment(ctx context.Context, header http.Header, body *circletype.DeleteCommentReq) (*circletype.DeleteCommentRes, error) {
	return Post[circletype.DeleteCommentRes](ctx, commonconfig.SocialUrl+"/internal/Circle.DeleteComment", header, body)
}

// RemoveCircleMembers 移除圈子成员（管理后台）。
func RemoveCircleMembers(ctx context.Context, header http.Header, body *circletype.RemoveCircleMembersReq) (*circletype.RemoveCircleMembersRes, error) {
	return Post[circletype.RemoveCircleMembersRes](ctx, commonconfig.SocialUrl+"/internal/Circle.RemoveCircleMembers", header, body)
}
