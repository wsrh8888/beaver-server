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
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateModerationCaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateModerationCaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModerationCaseLogic {
	return &CreateModerationCaseLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateModerationCaseLogic) CreateModerationCase(req *types.CreateModerationCaseReq) (resp *types.CreateModerationCaseRes, err error) {
	if req.TargetType <= 0 || req.TargetID == "" {
		return nil, errors.New("处置对象不能为空")
	}
	if req.Title == "" {
		return nil, errors.New("工单标题不能为空")
	}

	priority := req.Priority
	if priority <= 0 {
		priority = 1
	}

	caseRecord := backend_models.AdminModerationCase{
		CaseNo:      newCaseNo(),
		Source:      backend_models.CaseSourceManual,
		TargetType:  req.TargetType,
		TargetID:    req.TargetID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Status:      backend_models.CaseStatusPending,
	}
	if err = l.svcCtx.DB.Create(&caseRecord).Error; err != nil {
		l.Errorf("创建工单失败: %v", err)
		return nil, err
	}

	l.svcCtx.RecordOperation(req.UserID, "create_case", "case", fmt.Sprintf("%d", caseRecord.Id), uint64(caseRecord.Id),
		fmt.Sprintf("手动创建工单 targetType=%d targetId=%s", req.TargetType, req.TargetID), "success", "")

	return &types.CreateModerationCaseRes{CaseID: uint64(caseRecord.Id), CaseNo: caseRecord.CaseNo}, nil
}
