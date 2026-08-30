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
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendBlockListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// GetFriendBlockList 管理后台：好友黑名单列表。
// admin 职责：运营筛选条件映射、UserRpc 批量组装双方昵称。
// RPC 职责：ListFriendBlocks 领域查询，不与本 HTTP 接口 1:1。
func NewGetFriendBlockListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendBlockListLogic {
	return &GetFriendBlockListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendBlockListLogic) GetFriendBlockList(req *types.GetFriendBlockListReq) (resp *types.GetFriendBlockListRes, err error) {
	rpcRes, err := l.svcCtx.FriendRpc.ListFriendBlocks(l.ctx, &friend_rpc.ListFriendBlocksReq{
		Page:          int32(req.Page),
		PageSize:      int32(req.PageSize),
		UserId:        req.UserID,
		BlockedUserId: req.BlockedUserID,
	})
	if err != nil {
		l.Errorf("获取黑名单列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List)*2)
	for _, item := range rpcRes.List {
		for _, uid := range []string{item.UserId, item.BlockedUserId} {
			if uid == "" {
				continue
			}
			if _, ok := seen[uid]; ok {
				continue
			}
			seen[uid] = struct{}{}
			userIDs = append(userIDs, uid)
		}
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.GetFriendBlockListItem, 0, len(rpcRes.List))
	for _, item := range rpcRes.List {
		userName := ""
		if u, ok := users[item.UserId]; ok && u != nil {
			userName = u.NickName
		}
		blockedName := ""
		if u, ok := users[item.BlockedUserId]; ok && u != nil {
			blockedName = u.NickName
		}
		list = append(list, types.GetFriendBlockListItem{
			Id:              item.BlockId,
			UserID:          item.UserId,
			UserName:        userName,
			BlockedUserID:   item.BlockedUserId,
			BlockedUserName: blockedName,
			CreateTime:      item.CreatedAt,
		})
	}

	return &types.GetFriendBlockListRes{List: list, Total: rpcRes.Total}, nil
}
