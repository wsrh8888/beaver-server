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

type UpdateArchitectureLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：更新应用架构
func NewUpdateArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateArchitectureLogic {
	return &UpdateArchitectureLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("UpdateArchitecture", ctx),
	}
}

// UpdateArchitecture 只做协议转换，业务规则全部留在 platform_rpc。
func (l *UpdateArchitectureLogic) UpdateArchitecture(req *types.UpdateArchitectureReq) (resp *types.UpdateArchitectureRes, err error) {
	rpcReq := &platformcli.UpdateArchitectureReq{Id: req.ID, Description: req.Description, IsActive: req.IsActive}
	_, err = l.svcCtx.PlatformRpc.UpdateArchitecture(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.UpdateArchitectureRes{}, nil
}
