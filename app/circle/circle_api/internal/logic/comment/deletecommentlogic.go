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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type DeleteCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("delete_comment", ctx),
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentReq) (resp *types.DeleteCommentRes, err error) {
	var c circle_models.CircleCommentModel
	if err = l.svcCtx.DB.Where("comment_id = ? AND is_deleted = false", req.CommentID).First(&c).Error; err != nil {
		return nil, fmt.Errorf("评论不存在")
	}

	// 评论人可删，圈主/管理员也可删
	if c.UserID != req.UserID {
		var member circle_models.CircleMemberModel
		if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", c.CircleID, req.UserID).First(&member).Error; err != nil {
			return nil, fmt.Errorf("无权限删除")
		}
		if member.Role > 2 {
			return nil, fmt.Errorf("无权限删除")
		}
	}

	if err = l.svcCtx.DB.Model(&c).Update("is_deleted", true).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除评论失败", Data: map[string]any{"commentId": req.CommentID, "err": err.Error()}})
		return nil, fmt.Errorf("删除评论失败: %v", err)
	}

	l.logger.Info(model.LogMsg{Text: "删除评论成功", Data: map[string]interface{}{"commentId": req.CommentID, "userId": req.UserID}})

	return &types.DeleteCommentRes{}, nil
}
