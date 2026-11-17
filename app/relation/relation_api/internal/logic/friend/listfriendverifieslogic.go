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

type ListFriendVerifiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询好友验证记录
func NewListFriendVerifiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendVerifiesLogic {
	return &ListFriendVerifiesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListFriendVerifies", ctx),
	}
}

// ListFriendVerifies 只做协议转换，业务规则全部留在 relation_rpc。
func (l *ListFriendVerifiesLogic) ListFriendVerifies(req *types.ListFriendVerifiesReq) (resp *types.ListFriendVerifiesRes, err error) {
	rpcReq := &relationtypes.ListFriendVerifiesReq{Page: req.Page, PageSize: req.PageSize, SendUserId: req.SendUserID, RevUserId: req.RevUserID, SendStatus: req.SendStatus, RevStatus: req.RevStatus, StartTime: req.StartTime, EndTime: req.EndTime, VerifyId: req.VerifyID}
	rpcRes, err := l.svcCtx.FriendRpc.ListFriendVerifies(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.FriendVerifyItem, 0, len(rpcRes.List))
	for _, it1 := range rpcRes.List {
		tmp2 = append(tmp2, types.FriendVerifyItem{VerifyID: it1.VerifyId, SendUserID: it1.SendUserId, RevUserID: it1.RevUserId, SendStatus: it1.SendStatus, RevStatus: it1.RevStatus, Message: it1.Message, Source: it1.Source, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.ListFriendVerifiesRes{Total: rpcRes.Total, List: tmp2}, nil
}
