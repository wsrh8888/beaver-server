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

const friendActionRestore int32 = 2

type RestoreFriendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRestoreFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestoreFriendLogic {
	return &RestoreFriendLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// RestoreFriend 管理后台：恢复软删除的好友关系。
// admin 职责：校验 friendId，映射为 UpdateFriends 恢复 action。
// RPC 职责：UpdateFriends 统一处理关系状态变更。
func (l *RestoreFriendLogic) RestoreFriend(req *types.RestoreFriendReq) (resp *types.RestoreFriendRes, err error) {
	if req.FriendID == "" {
		return nil, errors.New("好友关系ID不能为空")
	}

	_, err = l.svcCtx.FriendRpc.UpdateFriends(l.ctx, &friend_rpc.UpdateFriendsReq{
		RelationIds: []string{req.FriendID},
		Action:      friendActionRestore,
	})
	if err != nil {
		l.Errorf("恢复好友失败: %v", err)
		return nil, err
	}
	return &types.RestoreFriendRes{}, nil
}
