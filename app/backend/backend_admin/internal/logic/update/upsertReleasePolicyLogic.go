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

// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpsertReleasePolicyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpsertReleasePolicyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpsertReleasePolicyLogic {
	return &UpsertReleasePolicyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpsertReleasePolicyLogic) UpsertReleasePolicy(req *types.UpsertReleasePolicyReq) (resp *types.UpsertReleasePolicyRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.UpsertReleasePolicy(l.ctx, &platform_rpc.UpsertReleasePolicyReq{
		AppId:           req.AppID,
		ArchitectureId:  uint64(req.ArchitectureID),
		StableVersionId: uint64(req.StableVersionID),
		GrayVersionId:   uint64(req.GrayVersionID),
		RolloutPercent:  uint32(req.RolloutPercent),
		MinVersion:      req.MinVersion,
		ForceUpdate:     req.ForceUpdate,
		IsActive:        req.IsActive,
	})
	if err != nil {
		l.Errorf("保存发版策略失败: %v", err)
		return nil, err
	}
	return &types.UpsertReleasePolicyRes{ID: uint(rpcRes.Id)}, nil
}
