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

type ListFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFriendsLogic {
	return &ListFriendsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListFriendsLogic) ListFriends(in *friend_rpc.ListFriendsReq) (*friend_rpc.ListFriendsRes, error) {
	if in.RelationId != "" {
		f, err := findFriend(l.svcCtx.DB, in.RelationId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &friend_rpc.ListFriendsRes{}, nil
		}
		if err != nil {
			return nil, err
		}
		return &friend_rpc.ListFriendsRes{
			Total: 1,
			List:  []*friend_rpc.FriendItem{toFriendItem(*f)},
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

	db := l.svcCtx.DB.Model(&friend_models.FriendModel{})
	if in.UserId != "" {
		db = db.Where("send_user_id = ? OR rev_user_id = ?", in.UserId, in.UserId)
	}
	if in.PeerUserId != "" {
		db = db.Where("send_user_id = ? OR rev_user_id = ?", in.PeerUserId, in.PeerUserId)
	}
	if in.IsDeleted != nil {
		db = db.Where("is_deleted = ?", *in.IsDeleted)
	} else {
		db = db.Where("is_deleted = false")
	}
	if in.StartTime != "" {
		db = db.Where("created_at >= ?", in.StartTime)
	}
	if in.EndTime != "" {
		db = db.Where("created_at <= ?", in.EndTime)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计好友失败: %v", err)
		return nil, err
	}

	var list []friend_models.FriendModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询好友列表失败: %v", err)
		return nil, err
	}

	items := make([]*friend_rpc.FriendItem, 0, len(list))
	for _, f := range list {
		items = append(items, toFriendItem(f))
	}
	return &friend_rpc.ListFriendsRes{Total: total, List: items}, nil
}
