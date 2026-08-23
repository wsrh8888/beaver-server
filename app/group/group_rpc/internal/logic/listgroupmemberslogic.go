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

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupMembersLogic {
	return &ListGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListGroupMembersLogic) ListGroupMembers(in *group_rpc.ListGroupMembersReq) (*group_rpc.ListGroupMembersRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&group_models.GroupMemberModel{})
	if in.GroupId != "" {
		db = db.Where("group_id = ?", in.GroupId)
	}
	if in.Role != 0 {
		db = db.Where("role = ?", in.Role)
	}
	if in.Status != 0 {
		db = db.Where("status = ?", in.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计群成员失败: %v", err)
		return nil, err
	}

	var list []group_models.GroupMemberModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询群成员列表失败: %v", err)
		return nil, err
	}

	items := make([]*group_rpc.GroupMemberItem, 0, len(list))
	for _, m := range list {
		items = append(items, &group_rpc.GroupMemberItem{
			Id:                 uint64(m.Id),
			GroupId:            m.GroupID,
			UserId:             m.UserID,
			Role:               int32(m.Role),
			Status:             int32(m.Status),
			ProhibitionMinutes: remainingMuteMinutes(m.MutedUntil),
			CreatedAt:          time.Time(m.CreatedAt).Format(time.RFC3339),
			UpdatedAt:          time.Time(m.UpdatedAt).Format(time.RFC3339),
		})
	}
	return &group_rpc.ListGroupMembersRes{Total: total, List: items}, nil
}

func remainingMuteMinutes(mutedUntil *time.Time) int32 {
	if mutedUntil == nil {
		return 0
	}
	remaining := time.Until(*mutedUntil)
	if remaining <= 0 {
		return 0
	}
	return int32(remaining.Minutes()) + 1
}
