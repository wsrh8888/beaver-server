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

package circle

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleCommentListLogic {
	return &GetCircleCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleCommentListLogic) GetCircleCommentList(req *types.GetCircleCommentListReq) (resp *types.GetCircleCommentListRes, err error) {
	if req.PostId == "" {
		return nil, errors.New("帖子ID不能为空")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	rpcRes, err := l.svcCtx.CircleRpc.ListComments(l.ctx, &circle_rpc.ListCommentsReq{
		PostId:   req.PostId,
		Page:     int32(page),
		PageSize: int32(limit),
	})
	if err != nil {
		l.Errorf("获取帖子评论列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetCircleCommentListItem, 0, len(rpcRes.List))
	for _, c := range rpcRes.List {
		list = append(list, types.GetCircleCommentListItem{
			CommentId: c.CommentId,
			PostId:    c.PostId,
			UserId:    c.UserId,
			Content:   c.Content,
			IsDeleted: c.IsDeleted,
			CreatedAt: c.CreatedAt,
		})
	}
	return &types.GetCircleCommentListRes{List: list, Total: rpcRes.Total}, nil
}
