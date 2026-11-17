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

type AddGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：添加群成员
func NewAddGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGroupMemberLogic {
	return &AddGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("AddGroupMember", ctx),
	}
}

// AddGroupMember 只做协议转换，业务规则全部留在 relation_rpc。
func (l *AddGroupMemberLogic) AddGroupMember(req *types.AddGroupMemberReq) (resp *types.AddGroupMemberRes, err error) {
	rpcReq := &relationtypes.AddGroupMemberReq{GroupId: req.GroupID, UserId: req.UserID, OperatedBy: req.OperatedBy}
	rpcRes, err := l.svcCtx.GroupRpc.AddGroupMember(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.AddGroupMemberRes{Added: rpcRes.Added, MemberVersion: rpcRes.MemberVersion}, nil
}
