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
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperationLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperationLogListLogic {
	return &GetOperationLogListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetOperationLogListLogic) GetOperationLogList(req *types.GetOperationLogListReq) (resp *types.GetOperationLogListRes, err error) {
	page, pageSize := req.Page, req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	db := l.svcCtx.DB.Model(&backend_models.AdminOperationLog{})
	if req.OperatorID != "" {
		db = db.Where("operator_id = ?", req.OperatorID)
	}
	if req.TargetID != "" {
		db = db.Where("target_id = ?", req.TargetID)
	}
	if req.TargetType != "" {
		db = db.Where("target_type = ?", req.TargetType)
	}
	if req.Actions != "" {
		parts := strings.Split(req.Actions, ",")
		db = db.Where("action IN ?", parts)
	} else if req.Action != "" {
		db = db.Where("action = ?", req.Action)
	}
	if req.CaseID > 0 {
		db = db.Where("case_id = ?", req.CaseID)
	}

	var total int64
	if err = db.Count(&total).Error; err != nil {
		l.Errorf("统计审计日志失败: %v", err)
		return nil, err
	}

	var rows []backend_models.AdminOperationLog
	if err = db.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		l.Errorf("查询审计日志失败: %v", err)
		return nil, err
	}

	list := make([]types.OperationLogInfo, 0, len(rows))
	for _, row := range rows {
		list = append(list, types.OperationLogInfo{
			ID:           uint64(row.Id),
			OperatorID:   row.OperatorID,
			Action:       row.Action,
			TargetType:   row.TargetType,
			TargetID:     row.TargetID,
			CaseID:       row.CaseID,
			Detail:       row.Detail,
			Result:       row.Result,
			ErrorMessage: row.ErrorMessage,
			CreatedAt:    row.CreatedAt.String(),
		})
	}
	return &types.GetOperationLogListRes{List: list, Total: total}, nil
}
