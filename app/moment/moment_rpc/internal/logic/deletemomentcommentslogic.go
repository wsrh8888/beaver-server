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

	"beaver/app/moment/moment_models"
	"beaver/app/moment/moment_rpc/internal/svc"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMomentCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMomentCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMomentCommentsLogic {
	return &DeleteMomentCommentsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteMomentCommentsLogic) DeleteMomentComments(in *moment_rpc.DeleteMomentCommentsReq) (*moment_rpc.DeleteMomentCommentsRes, error) {
	if len(in.CommentIds) == 0 {
		return &moment_rpc.DeleteMomentCommentsRes{}, nil
	}
	if err := l.svcCtx.DB.Where("comment_id IN ?", in.CommentIds).Delete(&moment_models.MomentCommentModel{}).Error; err != nil {
		l.Errorf("删除评论失败: %v", err)
		return nil, err
	}
	return &moment_rpc.DeleteMomentCommentsRes{}, nil
}
