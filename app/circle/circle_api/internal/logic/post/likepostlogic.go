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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type LikePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewLikePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikePostLogic {
	return &LikePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("like_post", ctx),
	}
}

func (l *LikePostLogic) LikePost(req *types.LikePostReq) (resp *types.LikePostRes, err error) {
	var p circle_models.CirclePostModel
	if err = l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", req.PostID).First(&p).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	var like circle_models.CircleLikeModel
	exists := l.svcCtx.DB.Where("post_id = ? AND user_id = ?", req.PostID, req.UserID).First(&like).Error == nil

	if req.Status && !exists {
		newLike := circle_models.CircleLikeModel{
			PostID:   req.PostID,
			UserID:   req.UserID,
			CircleID: p.CircleID,
		}
		if err = l.svcCtx.DB.Create(&newLike).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "点赞失败", Data: map[string]any{"postId": req.PostID, "userId": req.UserID, "err": err.Error()}})
			return nil, fmt.Errorf("点赞失败: %v", err)
		}
	} else if !req.Status && exists {
		if err = l.svcCtx.DB.Delete(&like).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "取消点赞失败", Data: map[string]any{"postId": req.PostID, "userId": req.UserID, "err": err.Error()}})
			return nil, fmt.Errorf("取消点赞失败: %v", err)
		}
	}

	l.logger.Info(model.LogMsg{Text: "点赞操作成功", Data: map[string]interface{}{"postId": req.PostID, "userId": req.UserID, "status": req.Status}})

	return &types.LikePostRes{}, nil
}
