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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupListLogic {
	return &GetGroupListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupListLogic) GetGroupList(req *types.GetGroupListReq) (resp *types.GetGroupListRes, err error) {
	rpcRes, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		Status:   int32(req.Status),
		Type:     int32(req.Type),
		Keywords: req.Keywords,
	})
	if err != nil {
		l.Errorf("获取群组列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetGroupListItem, 0, len(rpcRes.List))
	for _, g := range rpcRes.List {
		list = append(list, types.GetGroupListItem{
			Id:        uint(g.Id),
			GroupId:   g.GroupId,
			Type:      int(g.Type),
			Title:     g.Title,
			FileName:  g.Avatar,
			CreatorId: g.CreatorId,
			Notice:    g.Notice,
			Status:    int(g.Status),
			CreatedAt: g.CreatedAt,
			UpdatedAt: g.UpdatedAt,
		})
	}
	return &types.GetGroupListRes{List: list, Total: rpcRes.Total}, nil
}
