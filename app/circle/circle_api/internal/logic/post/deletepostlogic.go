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

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("delete_post", ctx),
	}
}

func (l *DeletePostLogic) DeletePost(req *types.DeletePostReq) (resp *types.DeletePostRes, err error) {
	var p circle_models.CirclePostModel
	if err = l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", req.PostID).First(&p).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	// 发帖人可删除，圈主/管理员也可删除
	if p.UserID != req.UserID {
		var member circle_models.CircleMemberModel
		if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", p.CircleID, req.UserID).First(&member).Error; err != nil {
			return nil, fmt.Errorf("无权限删除")
		}
		if member.Role > 2 {
			return nil, fmt.Errorf("无权限删除")
		}
	}

	if err = l.svcCtx.DB.Model(&p).Update("is_deleted", true).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "删除帖子失败", Data: map[string]any{"postId": req.PostID, "err": err.Error()}})
		return nil, fmt.Errorf("删除帖子失败: %v", err)
	}

	l.logger.Info(model.LogMsg{Text: "删除帖子成功", Data: map[string]interface{}{"postId": req.PostID, "userId": req.UserID}})

	return &types.DeletePostRes{}, nil
}
