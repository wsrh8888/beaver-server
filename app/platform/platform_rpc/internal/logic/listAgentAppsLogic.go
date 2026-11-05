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

	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAgentAppsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAgentAppsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentAppsLogic {
	return &ListAgentAppsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAgentAppsLogic) ListAgentApps(in *platform_rpc.ListAgentAppsReq) (*platform_rpc.ListAgentAppsRes, error) {
	q := l.svcCtx.DB.Model(&agent_models.AgentApp{})
	if in.Scene != "" {
		q = q.Where("scene = ?", in.Scene)
	}
	if in.Status != 0 {
		q = q.Where("status = ?", int(in.Status))
	}
	var total int64
	q.Count(&total)
	page := int(in.Page)
	pageSize := int(in.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var rows []agent_models.AgentApp
	q.Order("sort ASC, id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows)
	out := make([]*platform_rpc.AgentAppItem, 0, len(rows))
	for i := range rows {
		out = append(out, toAgentAppItem(&rows[i]))
	}
	return &platform_rpc.ListAgentAppsRes{Total: total, Apps: out}, nil
}
