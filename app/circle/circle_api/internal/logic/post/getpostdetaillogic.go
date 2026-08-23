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
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostDetailLogic {
	return &GetPostDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostDetailLogic) GetPostDetail(req *types.GetPostDetailReq) (resp *types.GetPostDetailRes, err error) {
	var p circle_models.CirclePostModel
	if err = l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", req.PostID).First(&p).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	// 用户信息
	userName, avatar := "", ""
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{p.UserID}})
	if userResp != nil {
		if info := userResp.UserInfo[p.UserID]; info != nil {
			userName = info.NickName
			avatar = info.Avatar
		}
	}

	// 是否点赞
	isLiked := false
	var like circle_models.CircleLikeModel
	if l.svcCtx.DB.Where("post_id = ? AND user_id = ?", req.PostID, req.UserID).First(&like).Error == nil {
		isLiked = true
	}

	// 最新20条一级评论
	var comments []circle_models.CircleCommentModel
	l.svcCtx.DB.Where("post_id = ? AND parent_id = '' AND is_deleted = false", req.PostID).
		Order("created_at DESC").Limit(20).Find(&comments)

	commentItems := buildCommentItems(l.ctx, l.svcCtx, comments, req.PostID)

	var commentCount int64
	l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Where("post_id = ? AND is_deleted = false", req.PostID).
		Count(&commentCount)

	var likeCount int64
	l.svcCtx.DB.Model(&circle_models.CircleLikeModel{}).
		Where("post_id = ?", req.PostID).
		Count(&likeCount)

	var likes []circle_models.CircleLikeModel
	l.svcCtx.DB.Where("post_id = ?", req.PostID).Order("created_at DESC").Limit(100).Find(&likes)
	likeItems := make([]types.GetPostDetailLikeInfo, 0, len(likes))
	if len(likes) > 0 {
		likeUserIDs := make([]string, 0, len(likes))
		for _, item := range likes {
			likeUserIDs = append(likeUserIDs, item.UserID)
		}
		likeUserResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: likeUserIDs})
		for _, item := range likes {
			likeName, likeAvatar := "", ""
			if likeUserResp != nil {
				if info := likeUserResp.UserInfo[item.UserID]; info != nil {
					likeName = info.NickName
					likeAvatar = info.Avatar
				}
			}
			likeItems = append(likeItems, types.GetPostDetailLikeInfo{
				UserID:   item.UserID,
				UserName: likeName,
				Avatar:   likeAvatar,
			})
		}
	}

	resp = &types.GetPostDetailRes{
		PostID:       p.PostID,
		CircleID:     p.CircleID,
		UserID:       p.UserID,
		UserName:     userName,
		Avatar:       avatar,
		Content:      p.Content,
		CommentCount: commentCount,
		LikeCount:    likeCount,
		IsLiked:      isLiked,
		IsTop:        p.IsTop,
		Comments:     commentItems,
		Likes:        likeItems,
		CreatedAt:    p.CreatedAt.String(),
	}
	if p.Files != nil {
		for _, f := range *p.Files {
			resp.Files = append(resp.Files, types.GetPostDetailFileInfo{FileKey: f.FileKey, Type: f.Type})
		}
	}
	return resp, nil
}

func buildCommentItems(ctx context.Context, svcCtx *svc.ServiceContext, comments []circle_models.CircleCommentModel, postID string) []types.GetPostDetailCommentInfo {
	if len(comments) == 0 {
		return []types.GetPostDetailCommentInfo{}
	}

	commentIDs := make([]string, 0, len(comments))
	userIDSet := make(map[string]struct{})
	for _, c := range comments {
		commentIDs = append(commentIDs, c.CommentID)
		userIDSet[c.UserID] = struct{}{}
		if c.ReplyToUserID != "" {
			userIDSet[c.ReplyToUserID] = struct{}{}
		}
	}

	type countResult struct {
		ParentID string
		Count    int64
	}
	var counts []countResult
	svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Select("parent_id, count(*) as count").
		Where("parent_id IN ? AND is_deleted = false", commentIDs).
		Group("parent_id").
		Scan(&counts)
	childCountMap := make(map[string]int64)
	for _, cr := range counts {
		childCountMap[cr.ParentID] = cr.Count
	}

	childMap := make(map[string][]circle_models.CircleCommentModel)
	var children []circle_models.CircleCommentModel
	svcCtx.DB.Where("parent_id IN ? AND is_deleted = false", commentIDs).
		Order("created_at ASC").
		Find(&children)
	for _, ch := range children {
		childMap[ch.ParentID] = append(childMap[ch.ParentID], ch)
		userIDSet[ch.UserID] = struct{}{}
		if ch.ReplyToUserID != "" {
			userIDSet[ch.ReplyToUserID] = struct{}{}
		}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}
	userResp, _ := svcCtx.UserRpc.UserListInfo(ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	fillUser := func(item *types.GetPostDetailCommentInfo, c circle_models.CircleCommentModel) {
		if userResp == nil {
			return
		}
		if info := userResp.UserInfo[c.UserID]; info != nil {
			item.UserName = info.NickName
			item.Avatar = info.Avatar
		}
		if c.ReplyToUserID != "" {
			if info := userResp.UserInfo[c.ReplyToUserID]; info != nil {
				item.ReplyToUserName = info.NickName
			}
		}
	}

	items := make([]types.GetPostDetailCommentInfo, 0, len(comments))
	for _, c := range comments {
		childItems := make([]types.GetPostDetailCommentInfo, 0, len(childMap[c.CommentID]))
		for _, ch := range childMap[c.CommentID] {
			childItem := types.GetPostDetailCommentInfo{
				CommentID:        ch.CommentID,
				UserID:           ch.UserID,
				Content:          ch.Content,
				ParentID:         ch.ParentID,
				ReplyToCommentID: ch.ReplyToCommentID,
				ChildCount:       0,
				Children:         []types.GetPostDetailCommentInfo{},
				CreatedAt:        ch.CreatedAt.String(),
			}
			fillUser(&childItem, ch)
			childItems = append(childItems, childItem)
		}
		item := types.GetPostDetailCommentInfo{
			CommentID:        c.CommentID,
			UserID:           c.UserID,
			Content:          c.Content,
			ParentID:         c.ParentID,
			ReplyToCommentID: c.ReplyToCommentID,
			ChildCount:       childCountMap[c.CommentID],
			Children:         childItems,
			CreatedAt:        c.CreatedAt.String(),
		}
		fillUser(&item, c)
		items = append(items, item)
	}
	return items
}
