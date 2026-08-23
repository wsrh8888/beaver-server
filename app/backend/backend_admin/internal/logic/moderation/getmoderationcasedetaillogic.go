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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetModerationCaseDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetModerationCaseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModerationCaseDetailLogic {
	return &GetModerationCaseDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetModerationCaseDetailLogic) GetModerationCaseDetail(req *types.GetModerationCaseDetailReq) (resp *types.GetModerationCaseDetailRes, err error) {
	if req.CaseID == 0 {
		return nil, errors.New("工单ID不能为空")
	}

	var c backend_models.AdminModerationCase
	if err = l.svcCtx.DB.Where("id = ?", req.CaseID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		l.Errorf("查询工单详情失败: %v", err)
		return nil, err
	}
	return &types.GetModerationCaseDetailRes{Case: formatCaseInfo(c)}, nil
}
