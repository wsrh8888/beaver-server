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

package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddArchitectureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddArchitectureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddArchitectureLogic {
	return &AddArchitectureLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AddArchitectureLogic) AddArchitecture(req *types.AddArchitectureReq) (resp *types.AddArchitectureRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateArchitecture(l.ctx, &platform_rpc.CreateArchitectureReq{
		AppId:       req.AppID,
		PlatformId:  uint32(req.PlatformID),
		ArchId:      uint32(req.ArchID),
		Description: req.Description,
	})
	if err != nil {
		l.Errorf("创建架构失败: %v", err)
		return nil, err
	}
	return &types.AddArchitectureRes{Id: uint(rpcRes.Id)}, nil
}
