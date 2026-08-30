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

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAppsLogic {
	return &ListAppsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListAppsLogic) ListApps(in *platform_rpc.ListAppsReq) (*platform_rpc.ListAppsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&platform_models.UpdateApp{})
	if in.IsActive != nil && *in.IsActive {
		db = db.Where("is_active = ?", true)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计应用失败: %v", err)
		return nil, err
	}

	var apps []platform_models.UpdateApp
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&apps).Error; err != nil {
		l.Errorf("查询应用列表失败: %v", err)
		return nil, err
	}

	list := make([]*platform_rpc.AppItem, 0, len(apps))
	for _, app := range apps {
		list = append(list, &platform_rpc.AppItem{
			Id:          uint64(app.Id),
			AppId:       app.AppID,
			Name:        app.Name,
			Description: app.Description,
			IsActive:    app.IsActive,
			CreatedAt:   app.CreatedAt.String(),
			UpdatedAt:   app.UpdatedAt.String(),
		})
	}

	return &platform_rpc.ListAppsRes{Total: total, Apps: list}, nil
}
