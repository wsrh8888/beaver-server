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

type GetCirclePostListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCirclePostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCirclePostListLogic {
	return &GetCirclePostListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCirclePostListLogic) GetCirclePostList(req *types.GetCirclePostListReq) (resp *types.GetCirclePostListRes, err error) {
	if req.CircleId == "" {
		return nil, errors.New("圈子ID不能为空")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	rpcRes, err := l.svcCtx.CircleRpc.ListPosts(l.ctx, &circle_rpc.ListPostsReq{
		CircleId: req.CircleId,
		Page:     int32(page),
		PageSize: int32(limit),
	})
	if err != nil {
		l.Errorf("获取圈子帖子列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetCirclePostListItem, 0, len(rpcRes.List))
	for _, p := range rpcRes.List {
		list = append(list, types.GetCirclePostListItem{
			PostId:       p.PostId,
			CircleId:     p.CircleId,
			UserId:       p.UserId,
			Content:      p.Content,
			IsDeleted:    p.IsDeleted,
			CommentCount: p.CommentCount,
			LikeCount:    p.LikeCount,
			CreatedAt:    p.CreatedAt,
		})
	}
	return &types.GetCirclePostListRes{List: list, Total: rpcRes.Total}, nil
}
