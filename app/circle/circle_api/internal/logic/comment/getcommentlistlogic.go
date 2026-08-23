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

package comment

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentListLogic) GetCommentList(req *types.GetCommentListReq) (resp *types.GetCommentListRes, err error) {
	var total int64
	var comments []circle_models.CircleCommentModel

	query := l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Where("post_id = ? AND is_deleted = false", req.PostID)
	if req.ParentID != "" {
		query = query.Where("parent_id = ?", req.ParentID)
	} else {
		query = query.Where("parent_id = ''")
	}
	query.Count(&total)
	query.Order("created_at ASC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&comments)

	if len(comments) == 0 {
		return &types.GetCommentListRes{Count: total, List: []types.GetCommentListItem{}}, nil
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

	childCountMap := make(map[string]int64)
	childMap := make(map[string][]circle_models.CircleCommentModel)
	// 仅顶层评论需要附带 children；按 parentId 拉子评论时不再嵌套
	if req.ParentID == "" {
		type countResult struct {
			ParentID string
			Count    int64
		}
		var counts []countResult
		l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
			Select("parent_id, count(*) as count").
			Where("parent_id IN ? AND is_deleted = false", commentIDs).
			Group("parent_id").
			Scan(&counts)
		for _, cr := range counts {
			childCountMap[cr.ParentID] = cr.Count
		}

		var children []circle_models.CircleCommentModel
		l.svcCtx.DB.Where("parent_id IN ? AND is_deleted = false", commentIDs).
			Order("created_at ASC").
			Find(&children)
		for _, ch := range children {
			childMap[ch.ParentID] = append(childMap[ch.ParentID], ch)
			userIDSet[ch.UserID] = struct{}{}
			if ch.ReplyToUserID != "" {
				userIDSet[ch.ReplyToUserID] = struct{}{}
			}
		}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	buildItem := func(c circle_models.CircleCommentModel, childCount int64, children []types.GetCommentListItem) types.GetCommentListItem {
		item := types.GetCommentListItem{
			CommentID:        c.CommentID,
			UserID:           c.UserID,
			Content:          c.Content,
			ParentID:         c.ParentID,
			ReplyToCommentID: c.ReplyToCommentID,
			ChildCount:       childCount,
			Children:         children,
			CreatedAt:        c.CreatedAt.String(),
		}
		if children == nil {
			item.Children = []types.GetCommentListItem{}
		}
		if userResp != nil {
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
		return item
	}

	items := make([]types.GetCommentListItem, 0, len(comments))
	for _, c := range comments {
		childItems := make([]types.GetCommentListItem, 0, len(childMap[c.CommentID]))
		for _, ch := range childMap[c.CommentID] {
			childItems = append(childItems, buildItem(ch, 0, []types.GetCommentListItem{}))
		}
		items = append(items, buildItem(c, childCountMap[c.CommentID], childItems))
	}

	return &types.GetCommentListRes{Count: total, List: items}, nil
}
