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

package system

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuthorityMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuthorityMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorityMenusLogic {
	return &GetAuthorityMenusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAuthorityMenusLogic) GetAuthorityMenus(req *types.GetAuthorityMenusReq) (resp *types.GetAuthorityMenusRes, err error) {
	var rows []backend_models.AdminSystemAuthorityMenu
	if err = l.svcCtx.DB.Where("authority_id = ?", req.Id).Find(&rows).Error; err != nil {
		l.Errorf("查询角色菜单失败: %v", err)
		return nil, err
	}
	menuIds := make([]uint, 0, len(rows))
	for _, row := range rows {
		menuIds = append(menuIds, uint(row.MenuID))
	}
	return &types.GetAuthorityMenusRes{MenuIds: menuIds}, nil
}
