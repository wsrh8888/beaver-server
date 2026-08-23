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
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupDetailLogic {
	return &GetGroupDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupDetailLogic) GetGroupDetail(req *types.GetGroupDetailReq) (resp *types.GetGroupDetailRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("群组ID不能为空")
	}

	rpcRes, err := l.svcCtx.GroupRpc.ListGroups(l.ctx, &group_rpc.ListGroupsReq{
		Id:       uint64(req.Id),
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		l.Errorf("获取群组详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("群组不存在")
	}
	g := rpcRes.List[0]
	return &types.GetGroupDetailRes{
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
	}, nil
}
