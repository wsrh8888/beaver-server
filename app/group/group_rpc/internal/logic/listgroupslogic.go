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

type ListGroupsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupsLogic {
	return &ListGroupsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListGroupsLogic) ListGroups(in *group_rpc.ListGroupsReq) (*group_rpc.ListGroupsRes, error) {
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

	db := l.svcCtx.DB.Model(&group_models.GroupModel{})
	if in.Id != 0 {
		db = db.Where("id = ?", in.Id)
	}
	if in.GroupId != "" {
		db = db.Where("group_id = ?", in.GroupId)
	}
	if in.Status != 0 {
		db = db.Where("status = ?", in.Status)
	}
	if in.Type != 0 {
		db = db.Where("type = ?", in.Type)
	}
	if in.Keywords != "" {
		db = db.Where("title LIKE ?", "%"+in.Keywords+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计群组失败: %v", err)
		return nil, err
	}

	var list []group_models.GroupModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询群组列表失败: %v", err)
		return nil, err
	}

	items := make([]*group_rpc.GroupItem, 0, len(list))
	for _, g := range list {
		items = append(items, &group_rpc.GroupItem{
			Id:        uint64(g.Id),
			GroupId:   g.GroupID,
			Type:      int32(g.Type),
			Title:     g.Title,
			Avatar:    g.Avatar,
			CreatorId: g.CreatorID,
			Notice:    g.Notice,
			Status:    int32(g.Status),
			MuteAll:   g.IsMuteAll,
			CreatedAt: time.Time(g.CreatedAt).Format(time.RFC3339),
			UpdatedAt: time.Time(g.UpdatedAt).Format(time.RFC3339),
		})
	}
	return &group_rpc.ListGroupsRes{Total: total, List: items}, nil
}
