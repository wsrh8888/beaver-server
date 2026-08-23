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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDevelopersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDevelopersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDevelopersLogic {
	return &ListDevelopersLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListDevelopersLogic) ListDevelopers(in *open_rpc.ListDevelopersReq) (*open_rpc.ListDevelopersRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&open_models.OpenDeveloper{})
	if in.Status > 0 {
		db = db.Where("status = ?", in.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计开发者失败: %v", err)
		return nil, err
	}

	var list []open_models.OpenDeveloper
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询开发者列表失败: %v", err)
		return nil, err
	}

	items := make([]*open_rpc.DeveloperItem, 0, len(list))
	for _, dev := range list {
		items = append(items, toDeveloperItem(dev))
	}

	return &open_rpc.ListDevelopersRes{Total: total, List: items}, nil
}
