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

type GetAgentToolLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentToolLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentToolLogic {
	return &GetAgentToolLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetAgentToolLogic) GetAgentTool(in *platform_rpc.GetAgentToolReq) (*platform_rpc.GetAgentToolRes, error) {
	var row agent_models.AgentTool
	if err := l.svcCtx.DB.Where("tool_id = ?", in.ToolId).First(&row).Error; err != nil {
		return &platform_rpc.GetAgentToolRes{}, nil
	}
	return &platform_rpc.GetAgentToolRes{Tool: toAgentToolItem(&row)}, nil
}

func toAgentToolItem(row *agent_models.AgentTool) *platform_rpc.AgentToolItem {
	if row == nil {
		return nil
	}
	return &platform_rpc.AgentToolItem{
		Id:       uint64(row.Id),
		ToolId:   row.ToolID,
		Name:     row.Name,
		ToolType: row.ToolType,
		Config:   row.Config,
		Status:   int32(row.Status),
		Sort:     int32(row.Sort),
	}
}
