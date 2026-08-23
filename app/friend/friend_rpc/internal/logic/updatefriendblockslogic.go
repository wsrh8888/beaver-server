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

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateFriendBlocksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendBlocksLogic {
	return &UpdateFriendBlocksLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateFriendBlocksLogic) UpdateFriendBlocks(in *friend_rpc.UpdateFriendBlocksReq) (*friend_rpc.UpdateFriendBlocksRes, error) {
	if in.Action != 1 {
		return nil, status.Error(codes.InvalidArgument, "无效的操作类型")
	}
	if len(in.BlockIds) == 0 {
		return &friend_rpc.UpdateFriendBlocksRes{}, nil
	}

	result := l.svcCtx.DB.Where("block_id IN ?", in.BlockIds).Delete(&friend_models.FriendBlockModel{})
	if result.Error != nil {
		l.Errorf("解除黑名单失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "操作失败")
	}
	if result.RowsAffected == 0 {
		result = l.svcCtx.DB.Where("id IN ?", in.BlockIds).Delete(&friend_models.FriendBlockModel{})
		if result.Error != nil {
			l.Errorf("解除黑名单失败: %v", result.Error)
			return nil, status.Error(codes.Internal, "操作失败")
		}
	}

	return &friend_rpc.UpdateFriendBlocksRes{AffectedCount: result.RowsAffected}, nil
}
