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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleListLogic {
	return &GetCircleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleListLogic) GetCircleList(req *types.GetCircleListReq) (resp *types.GetCircleListRes, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	rpcRes, err := l.svcCtx.CircleRpc.GetCircleList(l.ctx, &circle_rpc.GetCircleListReq{
		Page:     int32(page),
		PageSize: int32(limit),
		UserId:   req.UserId,
		Keywords: req.Keywords,
		CircleId: req.CircleId,
	})
	if err != nil {
		l.Errorf("获取圈子列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetCircleListItem, 0, len(rpcRes.List))
	for _, c := range rpcRes.List {
		list = append(list, mapCircleItem(c))
	}
	return &types.GetCircleListRes{List: list, Total: rpcRes.Total}, nil
}

func mapCircleItem(c *circle_rpc.CircleItem) types.GetCircleListItem {
	if c == nil {
		return types.GetCircleListItem{}
	}
	return types.GetCircleListItem{
		CircleId:    c.CircleId,
		Name:        c.Name,
		Description: c.Description,
		Avatar:      c.Avatar,
		CreatorId:   c.CreatorId,
		JoinType:    int(c.JoinType),
		MemberCount: c.MemberCount,
		PostCount:   c.PostCount,
		IsDeleted:   c.IsDeleted,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
