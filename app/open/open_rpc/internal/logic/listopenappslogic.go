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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ListOpenAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOpenAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOpenAppsLogic {
	return &ListOpenAppsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListOpenAppsLogic) ListOpenApps(in *open_rpc.ListOpenAppsReq) (*open_rpc.ListOpenAppsRes, error) {
	if in.AppId != "" {
		var app open_models.OpenApp
		if err := l.svcCtx.DB.Where("app_id = ?", in.AppId).First(&app).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &open_rpc.ListOpenAppsRes{}, nil
			}
			return nil, err
		}
		return &open_rpc.ListOpenAppsRes{
			Total: 1,
			List:  []*open_rpc.OpenAppItem{toOpenAppItem(app)},
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

	db := l.svcCtx.DB.Model(&open_models.OpenApp{})
	if in.Keyword != "" {
		like := "%" + in.Keyword + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", like, like)
	}
	if in.OwnerUserId != "" {
		db = db.Where("owner_user_id = ?", in.OwnerUserId)
	}
	if in.Status > 0 {
		db = db.Where("status = ?", in.Status-1) // 1草稿 2已发布 3禁用 -> 库内 0/1/2
	}
	if in.AuditStatus > 0 {
		db = db.Where("audit_status = ?", in.AuditStatus-1) // 1待审 2通过 3拒绝 -> 库内 0/1/2
	}
	switch in.CapabilityType {
	case 1:
		db = db.Where("enable_robot = ?", 1)
	case 2:
		db = db.Where("enable_webhook = ?", 1)
	case 3:
		db = db.Where("enable_robot = ? OR enable_webhook = ?", 1, 1)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计应用失败: %v", err)
		return nil, err
	}

	var list []open_models.OpenApp
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询应用列表失败: %v", err)
		return nil, err
	}

	items := make([]*open_rpc.OpenAppItem, 0, len(list))
	for _, app := range list {
		items = append(items, toOpenAppItem(app))
	}
	return &open_rpc.ListOpenAppsRes{Total: total, List: items}, nil
}
