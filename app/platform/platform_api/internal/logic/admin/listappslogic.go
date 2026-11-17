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
	beaverlog "beaver/utils/beaverlog"
)

type ListAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询应用
func NewListAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAppsLogic {
	return &ListAppsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListApps", ctx),
	}
}

// ListApps 只做协议转换，业务规则全部留在 platform_rpc。
func (l *ListAppsLogic) ListApps(req *types.ListAppsReq) (resp *types.ListAppsRes, err error) {
	rpcReq := &platformcli.ListAppsReq{Page: req.Page, PageSize: req.PageSize, IsActive: req.IsActive}
	rpcRes, err := l.svcCtx.PlatformRpc.ListApps(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.AppItem, 0, len(rpcRes.Apps))
	for _, it1 := range rpcRes.Apps {
		tmp2 = append(tmp2, types.AppItem{ID: it1.Id, AppID: it1.AppId, Name: it1.Name, Description: it1.Description, IsActive: it1.IsActive, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.ListAppsRes{Total: rpcRes.Total, Apps: tmp2}, nil
}
