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
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const friendActionHardDelete int32 = 1

type DeleteFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// DeleteFriend 管理后台：强制删除单条好友关系。
// admin 职责：校验 relationId，映射为 UpdateFriends 物理删除 action。
// RPC 职责：UpdateFriends 统一处理删除/恢复，不与 HTTP 路由 1:1。
func (l *DeleteFriendLogic) DeleteFriend(req *types.DeleteFriendReq) (resp *types.DeleteFriendRes, err error) {
	if req.FriendID == "" {
		return nil, errors.New("好友关系ID不能为空")
	}

	_, err = l.svcCtx.FriendRpc.UpdateFriends(l.ctx, &friend_rpc.UpdateFriendsReq{
		RelationIds: []string{req.FriendID},
		Action:      friendActionHardDelete,
	})
	if err != nil {
		l.Errorf("删除好友失败: %v", err)
		return nil, err
	}
	return &types.DeleteFriendRes{}, nil
}
