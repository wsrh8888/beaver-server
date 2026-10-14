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

package digital

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type ListDigitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDigitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalLogic {
	return &ListDigitalLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_digital", ctx)}
}

func (l *ListDigitalLogic) ListDigital(req *types.ListDigitalReq) (resp *types.ListDigitalRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	query := l.svcCtx.DB.Model(&open_models.OpenDigital{}).Where("owner_id = ?", req.UserID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	var partners []open_models.OpenDigital
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&partners).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	list := make([]types.DigitalInfo, 0, len(partners))
	for i := range partners {
		list = append(list, fillDigitalInfo(l.ctx, l.svcCtx, &partners[i]))
	}
	return &types.ListDigitalRes{Total: total, List: list}, nil
}
