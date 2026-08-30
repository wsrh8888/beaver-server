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

package workbench

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkbenchAppLogic {
	return &GetWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetWorkbenchAppLogic) GetWorkbenchApp(req *types.GetWorkbenchAppReq) (*types.GetWorkbenchAppRes, error) {
	rpcRes, err := l.svcCtx.PlatformRpc.GetWorkbenchApp(l.ctx, &platform_rpc.GetWorkbenchAppReq{
		WorkbenchAppId: req.WorkbenchAppID,
	})
	if err != nil {
		l.Errorf("获取工作台应用详情失败: %v", err)
		return nil, err
	}

	item := toAdminAppItem(rpcRes.App)
	return &types.GetWorkbenchAppRes{
		WorkbenchAppID: item.WorkbenchAppID,
		Name:           item.Name,
		Description:    item.Description,
		Icon:           item.Icon,
		AppType:        item.AppType,
		ClientScope:    item.ClientScope,
		EntryConfig:    item.EntryConfig,
		OpenMode:       item.OpenMode,
		Category:       item.Category,
		Sort:           item.Sort,
		Status:         item.Status,
		Remark:         item.Remark,
		CreatedBy:      item.CreatedBy,
		LastModifiedBy: item.LastModifiedBy,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}, nil
}
