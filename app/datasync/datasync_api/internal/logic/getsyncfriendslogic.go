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
	"time"

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/friend/friend_rpc/types/friend_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取所有需要更新的好友版本
func NewGetSyncFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncFriendsLogic {
	return &GetSyncFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_friends", ctx),
	}
}

func (l *GetSyncFriendsLogic) GetSyncFriends(req *types.GetSyncFriendsReq) (resp *types.GetSyncFriendsRes, err error) {
	// 调用Friend RPC获取好友版本信息
	friendResp, err := l.svcCtx.FriendRpc.GetFriendVersions(l.ctx, &friend_rpc.GetFriendVersionsReq{
		UserId: req.UserID,
		Since:  req.Since,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取好友版本信息失败", Data: map[string]any{"userId": req.UserID, "since": req.Since, "err": err.Error()}})
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "查询好友版本信息完成",
		Data: map[string]interface{}{"userId": req.UserID, "count": len(friendResp.FriendVersions)},
	})

	// 转换为响应格式，确保返回空数组而不是null
	friendVersions := make([]types.FriendVersionItem, 0)
	if friendResp.FriendVersions != nil {
		for _, friend := range friendResp.FriendVersions {
			friendVersions = append(friendVersions, types.FriendVersionItem{
				FriendId: friend.FriendId,
				Version:  friend.Version,
			})
		}
	}

	return &types.GetSyncFriendsRes{
		FriendVersions:  friendVersions,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
