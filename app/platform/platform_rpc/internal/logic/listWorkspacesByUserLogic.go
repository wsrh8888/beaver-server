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

	ai_models "beaver/app/platform/platform_models/ai"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWorkspacesByUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWorkspacesByUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWorkspacesByUserLogic {
	return &ListWorkspacesByUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListWorkspacesByUserLogic) ListWorkspacesByUser(in *platform_rpc.ListWorkspacesByUserReq) (*platform_rpc.ListWorkspacesByUserRes, error) {
	var rows []ai_models.AgentWorkspace
	if err := l.svcCtx.DB.Where("user_id = ?", in.UserId).
		Order("version DESC").Find(&rows).Error; err != nil {
		return &platform_rpc.ListWorkspacesByUserRes{}, nil
	}
	out := make([]*platform_rpc.GetLatestWorkspaceRes, 0, len(rows))
	for _, row := range rows {
		out = append(out, &platform_rpc.GetLatestWorkspaceRes{
			Id:       row.ID,
			UserId:   row.UserID,
			DeviceId: row.DeviceID,
			Name:     row.Name,
			Path:     row.Path,
			Version:  row.Version,
		})
	}
	return &platform_rpc.ListWorkspacesByUserRes{Workspaces: out}, nil
}
