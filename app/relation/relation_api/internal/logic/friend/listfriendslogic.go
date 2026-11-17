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

type ListFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询好友关系
func NewListFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendsLogic {
	return &ListFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListFriends", ctx),
	}
}

// ListFriends 只做协议转换，业务规则全部留在 relation_rpc。
func (l *ListFriendsLogic) ListFriends(req *types.ListFriendsReq) (resp *types.ListFriendsRes, err error) {
	rpcReq := &relationtypes.ListFriendsReq{Page: req.Page, PageSize: req.PageSize, UserId: req.UserID, PeerUserId: req.PeerUserID, IsDeleted: req.IsDeleted, StartTime: req.StartTime, EndTime: req.EndTime, RelationId: req.RelationID}
	rpcRes, err := l.svcCtx.FriendRpc.ListFriends(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.FriendItem, 0, len(rpcRes.List))
	for _, it1 := range rpcRes.List {
		tmp2 = append(tmp2, types.FriendItem{FriendID: it1.FriendId, SendUserID: it1.SendUserId, RevUserID: it1.RevUserId, SendUserNotice: it1.SendUserNotice, RevUserNotice: it1.RevUserNotice, Source: it1.Source, IsDeleted: it1.IsDeleted, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.ListFriendsRes{Total: rpcRes.Total, List: tmp2}, nil
}
