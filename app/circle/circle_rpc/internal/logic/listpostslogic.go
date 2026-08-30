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

package logic

import (
	"context"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPostsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPostsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPostsLogic {
	return &ListPostsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPostsLogic) ListPosts(in *circle_rpc.ListPostsReq) (*circle_rpc.ListPostsRes, error) {
	query := l.svcCtx.DB.Model(&circle_models.CirclePostModel{}).
		Where("circle_id = ? AND is_deleted = false", in.CircleId)

	var total int64
	query.Count(&total)

	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	var posts []circle_models.CirclePostModel
	l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", in.CircleId).
		Order("is_top DESC, created_at DESC").
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Find(&posts)

	postIDs := make([]string, 0, len(posts))
	for _, p := range posts {
		postIDs = append(postIDs, p.PostID)
	}

	commentCountMap := make(map[string]int64)
	likeCountMap := make(map[string]int64)
	if len(postIDs) > 0 {
		type countRow struct {
			PostID string
			Count  int64
		}
		var commentRows []countRow
		l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
			Select("post_id, count(*) as count").
			Where("post_id IN ? AND is_deleted = false", postIDs).
			Group("post_id").
			Scan(&commentRows)
		for _, row := range commentRows {
			commentCountMap[row.PostID] = row.Count
		}

		var likeRows []countRow
		l.svcCtx.DB.Model(&circle_models.CircleLikeModel{}).
			Select("post_id, count(*) as count").
			Where("post_id IN ?", postIDs).
			Group("post_id").
			Scan(&likeRows)
		for _, row := range likeRows {
			likeCountMap[row.PostID] = row.Count
		}
	}

	list := make([]*circle_rpc.PostItem, 0, len(posts))
	for _, p := range posts {
		list = append(list, &circle_rpc.PostItem{
			PostId:       p.PostID,
			CircleId:     p.CircleID,
			UserId:       p.UserID,
			Content:      p.Content,
			IsDeleted:    p.IsDeleted,
			CreatedAt:    p.CreatedAt.String(),
			CommentCount: commentCountMap[p.PostID],
			LikeCount:    likeCountMap[p.PostID],
		})
	}

	return &circle_rpc.ListPostsRes{Total: total, List: list}, nil
}
