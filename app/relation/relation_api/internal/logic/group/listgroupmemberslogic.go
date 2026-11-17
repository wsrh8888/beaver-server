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

type ListGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询群成员
func NewListGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupMembersLogic {
	return &ListGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListGroupMembers", ctx),
	}
}

// ListGroupMembers 只做协议转换，业务规则全部留在 relation_rpc。
func (l *ListGroupMembersLogic) ListGroupMembers(req *types.ListGroupMembersReq) (resp *types.ListGroupMembersRes, err error) {
	rpcReq := &relationtypes.ListGroupMembersReq{GroupId: req.GroupID, Page: req.Page, PageSize: req.PageSize, Role: req.Role, Status: req.Status}
	rpcRes, err := l.svcCtx.GroupRpc.ListGroupMembers(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.GroupMemberItem, 0, len(rpcRes.List))
	for _, it1 := range rpcRes.List {
		tmp2 = append(tmp2, types.GroupMemberItem{ID: it1.Id, GroupID: it1.GroupId, UserID: it1.UserId, Role: it1.Role, Status: it1.Status, ProhibitionMinutes: it1.ProhibitionMinutes, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.ListGroupMembersRes{Total: rpcRes.Total, List: tmp2}, nil
}
