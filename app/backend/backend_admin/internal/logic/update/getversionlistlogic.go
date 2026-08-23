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

type GetVersionListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVersionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVersionListLogic {
	return &GetVersionListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetVersionListLogic) GetVersionList(req *types.GetVersionListReq) (resp *types.GetVersionListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListVersions(l.ctx, &platform_rpc.ListVersionsReq{
		ArchitectureId: uint64(req.ArchitectureID),
		Page:           int32(req.Page),
		PageSize:       int32(req.PageSize),
	})
	if err != nil {
		l.Errorf("获取版本列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetVersionListItem, 0, len(rpcRes.Versions))
	for _, ver := range rpcRes.Versions {
		list = append(list, types.GetVersionListItem{
			VersionID:      uint(ver.VersionId),
			ArchitectureID: uint(ver.ArchitectureId),
			Version:        ver.Version,
			FileUrl:        ver.FileUrl,
			Description:    ver.Description,
			ReleaseNotes:   ver.ReleaseNotes,
			CreatedAt:      ver.CreatedAt,
			UpdatedAt:      ver.UpdatedAt,
		})
	}

	return &types.GetVersionListRes{Total: rpcRes.Total, Versions: list}, nil
}
