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

package moderation

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetModerationCaseListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetModerationCaseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModerationCaseListLogic {
	return &GetModerationCaseListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetModerationCaseListLogic) GetModerationCaseList(req *types.GetModerationCaseListReq) (resp *types.GetModerationCaseListRes, err error) {
	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	db := l.svcCtx.DB.Model(&backend_models.AdminModerationCase{})
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.TargetType > 0 {
		db = db.Where("target_type = ?", req.TargetType)
	}
	if req.Keyword != "" {
		kw := "%" + req.Keyword + "%"
		db = db.Where("case_no LIKE ? OR title LIKE ? OR target_id LIKE ?", kw, kw, kw)
	}

	var total int64
	if err = db.Count(&total).Error; err != nil {
		l.Errorf("统计工单失败: %v", err)
		return nil, err
	}

	var rows []backend_models.AdminModerationCase
	if err = db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.Errorf("查询工单列表失败: %v", err)
		return nil, err
	}

	list := make([]types.ModerationCaseInfo, 0, len(rows))
	for _, c := range rows {
		list = append(list, formatCaseInfo(c))
	}
	return &types.GetModerationCaseListRes{List: list, Total: total}, nil
}
