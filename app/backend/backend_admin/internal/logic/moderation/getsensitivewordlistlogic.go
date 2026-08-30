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

type GetSensitiveWordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSensitiveWordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSensitiveWordListLogic {
	return &GetSensitiveWordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSensitiveWordListLogic) GetSensitiveWordList(req *types.GetSensitiveWordListReq) (resp *types.GetSensitiveWordListRes, err error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.DB.Model(&backend_models.AdminSensitiveWord{})
	if req.Keyword != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("word LIKE ? OR category LIKE ?", like, like)
	}
	if req.IsActive {
		db = db.Where("is_active = ?", true)
	}

	var total int64
	if err = db.Count(&total).Error; err != nil {
		l.Errorf("统计敏感词失败: %v", err)
		return nil, err
	}

	var rows []backend_models.AdminSensitiveWord
	if err = db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.Errorf("查询敏感词失败: %v", err)
		return nil, err
	}

	list := make([]types.SensitiveWordInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, types.SensitiveWordInfo{
			ID:        uint64(row.Id),
			Word:      row.Word,
			Category:  row.Category,
			Level:     row.Level,
			IsActive:  row.IsActive,
			Remark:    row.Remark,
			CreatedAt: row.CreatedAt.String(),
		})
	}

	return &types.GetSensitiveWordListRes{List: list, Total: total}, nil
}
