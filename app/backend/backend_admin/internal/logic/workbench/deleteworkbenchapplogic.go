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

type DeleteWorkbenchAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkbenchAppLogic {
	return &DeleteWorkbenchAppLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteWorkbenchAppLogic) DeleteWorkbenchApp(req *types.DeleteWorkbenchAppReq) (*types.DeleteWorkbenchAppRes, error) {
	_, err := l.svcCtx.PlatformRpc.DeleteWorkbenchApp(l.ctx, &platform_rpc.DeleteWorkbenchAppReq{
		WorkbenchAppId: req.WorkbenchAppID,
	})
	if err != nil {
		l.Errorf("删除工作台应用失败: %v", err)
		return nil, err
	}

	return &types.DeleteWorkbenchAppRes{}, nil
}
