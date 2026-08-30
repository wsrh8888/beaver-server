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
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (resp *types.CreateCommentRes, err error) {
	var p circle_models.CirclePostModel
	if err = l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", req.PostID).First(&p).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	// 查被回复用户ID
	replyToUserID := ""
	if req.ReplyToCommentID != "" {
		var replyTo circle_models.CircleCommentModel
		if l.svcCtx.DB.Where("comment_id = ?", req.ReplyToCommentID).First(&replyTo).Error == nil {
			replyToUserID = replyTo.UserID
		}
	}

	commentID := uuid.New().String()
	c := circle_models.CircleCommentModel{
		CommentID:        commentID,
		PostID:           req.PostID,
		CircleID:         p.CircleID,
		UserID:           req.UserID,
		Content:          req.Content,
		ParentID:         req.ParentID,
		ReplyToCommentID: req.ReplyToCommentID,
		ReplyToUserID:    replyToUserID,
	}
	if err = l.svcCtx.DB.Create(&c).Error; err != nil {
		return nil, fmt.Errorf("发布评论失败: %v", err)
	}

	// 拉用户信息
	userName, avatar, replyToUserName := "", "", ""
	queryIDs := []string{req.UserID}
	if replyToUserID != "" {
		queryIDs = append(queryIDs, replyToUserID)
	}
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: queryIDs})
	if userResp != nil {
		if info := userResp.UserInfo[req.UserID]; info != nil {
			userName = info.NickName
			avatar = info.Avatar
		}
		if replyToUserID != "" {
			if info := userResp.UserInfo[replyToUserID]; info != nil {
				replyToUserName = info.NickName
			}
		}
	}

	return &types.CreateCommentRes{
		CommentID:        commentID,
		PostID:           req.PostID,
		UserID:           req.UserID,
		UserName:         userName,
		Avatar:           avatar,
		Content:          req.Content,
		ParentID:         req.ParentID,
		ReplyToCommentID: req.ReplyToCommentID,
		ReplyToUserName:  replyToUserName,
		CreatedAt:        c.CreatedAt.String(),
	}, nil
}
