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

type CreateWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkbenchAppLogic {
	return &CreateWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateWorkbenchAppLogic) CreateWorkbenchApp(req *types.CreateWorkbenchAppReq) (*types.CreateWorkbenchAppRes, error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateWorkbenchApp(l.ctx, &platform_rpc.CreateWorkbenchAppReq{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		AppType:     int32(req.AppType),
		ClientScope: int32(req.ClientScope),
		EntryConfig: toProtoEntryConfig(&req.EntryConfig),
		OpenMode:    int32(req.OpenMode),
		Category:    int32(req.Category),
		Sort:        int32(req.Sort),
		Status:      int32(req.Status),
		Remark:      req.Remark,
		OperatorId:  req.UserID,
	})
	if err != nil {
		l.Errorf("创建工作台应用失败: %v", err)
		return nil, err
	}

	return &types.CreateWorkbenchAppRes{WorkbenchAppID: rpcRes.WorkbenchAppId}, nil
}
