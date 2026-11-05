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

type GetAgentMcpServerLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentMcpServerLogic {
	return &GetAgentMcpServerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// —— MCP 服务目录（agent/ 通用）——
func (l *GetAgentMcpServerLogic) GetAgentMcpServer(in *platform_rpc.GetAgentMcpServerReq) (*platform_rpc.GetAgentMcpServerRes, error) {
	var row agent_models.AgentMcpServer
	if err := l.svcCtx.DB.Where("mcp_server_id = ?", in.McpServerId).First(&row).Error; err != nil {
		return &platform_rpc.GetAgentMcpServerRes{}, nil
	}
	return &platform_rpc.GetAgentMcpServerRes{Server: toAgentMcpServerItem(&row)}, nil
}

func toAgentMcpServerItem(row *agent_models.AgentMcpServer) *platform_rpc.AgentMcpServerItem {
	if row == nil {
		return nil
	}
	return &platform_rpc.AgentMcpServerItem{
		Id:           uint64(row.Id),
		McpServerId:  row.McpServerID,
		Name:         row.Name,
		ServerName:   row.ServerName,
		Transport:    row.Transport,
		Config:       row.Config,
		Tools:        row.Tools,
		Status:       int32(row.Status),
	}
}
