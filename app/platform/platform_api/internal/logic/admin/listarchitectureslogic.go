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

type ListArchitecturesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：分页查询应用架构
func NewListArchitecturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListArchitecturesLogic {
	return &ListArchitecturesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListArchitectures", ctx),
	}
}

// ListArchitectures 只做协议转换，业务规则全部留在 platform_rpc。
func (l *ListArchitecturesLogic) ListArchitectures(req *types.ListArchitecturesReq) (resp *types.ListArchitecturesRes, err error) {
	rpcReq := &platformcli.ListArchitecturesReq{AppId: req.AppID, Page: req.Page, PageSize: req.PageSize, IsActive: req.IsActive}
	rpcRes, err := l.svcCtx.PlatformRpc.ListArchitectures(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.ArchitectureItem, 0, len(rpcRes.Architectures))
	for _, it1 := range rpcRes.Architectures {
		tmp2 = append(tmp2, types.ArchitectureItem{ID: it1.Id, AppID: it1.AppId, AppName: it1.AppName, PlatformID: it1.PlatformId, ArchID: it1.ArchId, Description: it1.Description, IsActive: it1.IsActive, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.ListArchitecturesRes{Total: rpcRes.Total, Architectures: tmp2}, nil
}
