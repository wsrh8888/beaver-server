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

	"beaver/app/relation/relation_api/internal/svc"
	"beaver/app/relation/relation_api/internal/types"
	relationtypes "beaver/app/relation/relation_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
)

type GroupsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：按ID批量取群组（支持增量）
func NewGroupsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupsListByIdsLogic {
	return &GroupsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GroupsListByIds", ctx),
	}
}

// GroupsListByIds 只做协议转换，业务规则全部留在 relation_rpc。
func (l *GroupsListByIdsLogic) GroupsListByIds(req *types.GetGroupsListByIdsReq) (resp *types.GetGroupsListByIdsRes, err error) {
	rpcReq := &relationtypes.GetGroupsListByIdsReq{GroupIDs: req.GroupIDs, Since: req.Since}
	rpcRes, err := l.svcCtx.GroupRpc.GetGroupsListByIds(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.GroupListById, 0, len(rpcRes.Groups))
	for _, it1 := range rpcRes.Groups {
		tmp2 = append(tmp2, types.GroupListById{GroupID: it1.GroupID, Name: it1.Name, Avatar: it1.Avatar, Description: it1.Description, Version: it1.Version, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.GetGroupsListByIdsRes{Groups: tmp2}, nil
}
