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

type ListReleasePoliciesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：查询发布策略
func NewListReleasePoliciesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListReleasePoliciesLogic {
	return &ListReleasePoliciesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListReleasePolicies", ctx),
	}
}

// ListReleasePolicies 只做协议转换，业务规则全部留在 platform_rpc。
func (l *ListReleasePoliciesLogic) ListReleasePolicies(req *types.ListReleasePoliciesReq) (resp *types.ListReleasePoliciesRes, err error) {
	rpcReq := &platformcli.ListReleasePoliciesReq{AppId: req.AppID}
	rpcRes, err := l.svcCtx.PlatformRpc.ListReleasePolicies(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.ReleasePolicyItem, 0, len(rpcRes.Policies))
	for _, it1 := range rpcRes.Policies {
		tmp2 = append(tmp2, types.ReleasePolicyItem{ID: it1.Id, AppID: it1.AppId, ArchitectureID: it1.ArchitectureId, StableVersionID: it1.StableVersionId, GrayVersionID: it1.GrayVersionId, RolloutPercent: it1.RolloutPercent, MinVersion: it1.MinVersion, ForceUpdate: it1.ForceUpdate, IsActive: it1.IsActive, StableVersion: it1.StableVersion, GrayVersion: it1.GrayVersion, CreatedAt: it1.CreatedAt, UpdatedAt: it1.UpdatedAt})
	}
	return &types.ListReleasePoliciesRes{Policies: tmp2}, nil
}
