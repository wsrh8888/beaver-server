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

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ListFriendBlocksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFriendBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendBlocksLogic {
	return &ListFriendBlocksLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListFriendBlocksLogic) ListFriendBlocks(in *friend_rpc.ListFriendBlocksReq) (*friend_rpc.ListFriendBlocksRes, error) {
	if in.BlockId != "" {
		b, err := findFriendBlock(l.svcCtx.DB, in.BlockId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &friend_rpc.ListFriendBlocksRes{}, nil
		}
		if err != nil {
			return nil, err
		}
		return &friend_rpc.ListFriendBlocksRes{
			Total: 1,
			List:  []*friend_rpc.FriendBlockItem{toFriendBlockItem(*b)},
		}, nil
	}

	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&friend_models.FriendBlockModel{})
	if in.UserId != "" {
		db = db.Where("user_id = ?", in.UserId)
	}
	if in.BlockedUserId != "" {
		db = db.Where("blocked_user_id = ?", in.BlockedUserId)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计黑名单失败: %v", err)
		return nil, err
	}

	var list []friend_models.FriendBlockModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询黑名单列表失败: %v", err)
		return nil, err
	}

	items := make([]*friend_rpc.FriendBlockItem, 0, len(list))
	for _, b := range list {
		items = append(items, toFriendBlockItem(b))
	}
	return &friend_rpc.ListFriendBlocksRes{Total: total, List: items}, nil
}
