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

type GetArchitecturesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetArchitecturesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArchitecturesLogic {
	return &GetArchitecturesLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetArchitecturesLogic) GetArchitectures(req *types.GetArchitecturesReq) (resp *types.GetArchitecturesRes, err error) {
	rpcReq := &platform_rpc.ListArchitecturesReq{
		AppId:    req.AppID,
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
	}
	if req.IsActive {
		active := true
		rpcReq.IsActive = &active
	}

	rpcRes, err := l.svcCtx.PlatformRpc.ListArchitectures(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("获取架构列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetArchitecturesItem, 0, len(rpcRes.Architectures))
	for _, arch := range rpcRes.Architectures {
		list = append(list, types.GetArchitecturesItem{
			Id:          uint(arch.Id),
			AppID:       arch.AppId,
			AppName:     arch.AppName,
			PlatformID:  uint(arch.PlatformId),
			ArchID:      uint(arch.ArchId),
			Description: arch.Description,
			IsActive:    arch.IsActive,
			CreatedAt:   arch.CreatedAt,
			UpdatedAt:   arch.UpdatedAt,
		})
	}

	return &types.GetArchitecturesRes{Total: rpcRes.Total, Architectures: list}, nil
}
