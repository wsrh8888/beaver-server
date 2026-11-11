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

type ListAgentToolsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAgentToolsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentToolsLogic {
	return &ListAgentToolsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— 工具目录（agent/ 通用）——
func (l *ListAgentToolsLogic) ListAgentTools(in *platform_rpc.ListAgentToolsReq) (*platform_rpc.ListAgentToolsRes, error) {
	q := l.svcCtx.DB.Model(&agent_models.AgentTool{})
	if in.Status != 0 {
		q = q.Where("status = ?", int(in.Status))
	}
	if in.Source != "" {
		q = q.Where("source = ?", in.Source)
	}
	var rows []agent_models.AgentTool
	q.Order("sort ASC, id ASC").Find(&rows)
	out := make([]*platform_rpc.AgentToolItem, 0, len(rows))
	for i := range rows {
		out = append(out, toAgentToolItem(&rows[i]))
	}
	return &platform_rpc.ListAgentToolsRes{Tools: out}, nil
}
