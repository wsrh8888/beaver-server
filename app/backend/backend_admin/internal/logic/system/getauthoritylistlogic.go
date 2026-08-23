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

type GetAuthorityListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAuthorityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuthorityListLogic {
	return &GetAuthorityListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAuthorityListLogic) GetAuthorityList(req *types.GetAuthorityListReq) (resp *types.GetAuthorityListRes, err error) {
	var rows []backend_models.AdminSystemAuthority
	if err = l.svcCtx.DB.Order("sort ASC, id ASC").Find(&rows).Error; err != nil {
		l.Errorf("查询角色列表失败: %v", err)
		return nil, err
	}

	list := make([]types.AuthorityInfo, 0, len(rows))
	for _, row := range rows {
		var menuCount int64
		_ = l.svcCtx.DB.Model(&backend_models.AdminSystemAuthorityMenu{}).
			Where("authority_id = ?", row.Id).Count(&menuCount).Error
		list = append(list, types.AuthorityInfo{
			Id:          uint(row.Id),
			Name:        row.Name,
			Description: row.Description,
			Status:      row.Status,
			Sort:        row.Sort,
			MenuCount:   menuCount,
		})
	}
	return &types.GetAuthorityListRes{List: list}, nil
}
