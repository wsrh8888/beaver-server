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

package admin

import (
	"context"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	platformcli "beaver/app/platform/platform_rpc/platform"
	"beaver/common/convert"
	beaverlog "beaver/utils/beaverlog"
)

type AdminListWorkbenchAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询工作台应用
func NewAdminListWorkbenchAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListWorkbenchAppsLogic {
	return &AdminListWorkbenchAppsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("AdminListWorkbenchApps", ctx),
	}
}

// AdminListWorkbenchApps 只做协议转换，业务规则全部留在 platform_rpc。
// 对端 RPC 方法名是 ListWorkbenchApps（运营视角与 C 端共用同一实现，只是入参不同）。
func (l *AdminListWorkbenchAppsLogic) AdminListWorkbenchApps(req *types.AdminListWorkbenchAppsReq) (resp *types.AdminListWorkbenchAppsRes, err error) {
	rpcReq, err := convert.To[platformcli.ListWorkbenchAppsReq](req)
	if err != nil {
		return nil, err
	}
	rpcRes, err := l.svcCtx.PlatformRpc.ListWorkbenchApps(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return convert.To[types.AdminListWorkbenchAppsRes](rpcRes)
}
