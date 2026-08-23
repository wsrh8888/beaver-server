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

package post

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostListLogic {
	return &GetPostListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostListLogic) GetPostList(req *types.GetPostListReq) (resp *types.GetPostListRes, err error) {
	var total int64
	var posts []circle_models.CirclePostModel

	l.svcCtx.DB.Model(&circle_models.CirclePostModel{}).
		Where("circle_id = ? AND is_deleted = false", req.CircleID).
		Count(&total)
	l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", req.CircleID).
		Order("is_top DESC, created_at DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&posts)

	if len(posts) == 0 {
		return &types.GetPostListRes{Count: total, List: []types.PostListItem{}}, nil
	}

	userIDs := make([]string, 0, len(posts))
	postIDs := make([]string, 0, len(posts))
	for _, p := range posts {
		userIDs = append(userIDs, p.UserID)
		postIDs = append(postIDs, p.PostID)
	}
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	var likes []circle_models.CircleLikeModel
	l.svcCtx.DB.Where("post_id IN ? AND user_id = ?", postIDs, req.UserID).Find(&likes)
	likedMap := make(map[string]bool)
	for _, lk := range likes {
		likedMap[lk.PostID] = true
	}

	likeCountMap := make(map[string]int64)
	type likeCountRow struct {
		PostID string
		Count  int64
	}
	var likeRows []likeCountRow
	l.svcCtx.DB.Model(&circle_models.CircleLikeModel{}).
		Select("post_id, count(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id").
		Scan(&likeRows)
	for _, row := range likeRows {
		likeCountMap[row.PostID] = row.Count
	}

	commentCountMap := make(map[string]int64)
	type commentCountRow struct {
		PostID string
		Count  int64
	}
	var commentRows []commentCountRow
	l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Select("post_id, count(*) as count").
		Where("post_id IN ? AND is_deleted = false", postIDs).
		Group("post_id").
		Scan(&commentRows)
	for _, row := range commentRows {
		commentCountMap[row.PostID] = row.Count
	}

	commentMap := make(map[string][]circle_models.CircleCommentModel)
	var allComments []circle_models.CircleCommentModel
	l.svcCtx.DB.Where("post_id IN ? AND parent_id = '' AND is_deleted = false", postIDs).
		Order("created_at DESC").
		Find(&allComments)
	previewCountMap := make(map[string]int)
	for _, c := range allComments {
		if previewCountMap[c.PostID] < 3 {
			commentMap[c.PostID] = append(commentMap[c.PostID], c)
			previewCountMap[c.PostID]++
		}
	}

	items := make([]types.PostListItem, 0, len(posts))
	for _, p := range posts {
		item := types.PostListItem{
			PostID:       p.PostID,
			CircleID:     p.CircleID,
			UserID:       p.UserID,
			Content:      p.Content,
			CommentCount: commentCountMap[p.PostID],
			LikeCount:    likeCountMap[p.PostID],
			IsLiked:      likedMap[p.PostID],
			IsTop:        p.IsTop,
			CreatedAt:    p.CreatedAt.String(),
		}
		if userResp != nil {
			if info := userResp.UserInfo[p.UserID]; info != nil {
				item.UserName = info.NickName
				item.Avatar = info.Avatar
			}
		}
		if p.Files != nil {
			for _, f := range *p.Files {
				item.Files = append(item.Files, types.GetPostListFileInfo{FileKey: f.FileKey, Type: f.Type})
			}
		}
		for _, c := range commentMap[p.PostID] {
			item.Comments = append(item.Comments, types.GetPostListCommentInfo{
				CommentID: c.CommentID,
				UserID:    c.UserID,
				Content:   c.Content,
				CreatedAt: c.CreatedAt.String(),
			})
		}
		items = append(items, item)
	}

	return &types.GetPostListRes{Count: total, List: items}, nil
}
