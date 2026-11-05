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

type ListAgentMcpServersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAgentMcpServersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentMcpServersLogic {
	return &ListAgentMcpServersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// —— MCP 服务目录（agent/ 通用）——
func (l *ListAgentMcpServersLogic) ListAgentMcpServers(in *platform_rpc.ListAgentMcpServersReq) (*platform_rpc.ListAgentMcpServersRes, error) {
	q := l.svcCtx.DB.Model(&agent_models.AgentMcpServer{})
	if in.Status != 0 {
		q = q.Where("status = ?", int(in.Status))
	}
	if in.Source != "" {
		q = q.Where("source = ?", in.Source)
	}
	var rows []agent_models.AgentMcpServer
	q.Order("id ASC").Find(&rows)
	out := make([]*platform_rpc.AgentMcpServerItem, 0, len(rows))
	for _, row := range rows {
		out = append(out, &platform_rpc.AgentMcpServerItem{
			Id:         uint64(row.Id),
			McpServerId: row.McpServerID,
			Name:       row.Name,
			ServerName: row.ServerName,
			Transport:  row.Transport,
			Config:     row.Config,
			Tools:      row.Tools,
			Status:     int32(row.Status),
			Source:     row.Source,
			UserId:     row.UserID,
		})
	}
	return &platform_rpc.ListAgentMcpServersRes{Servers: out}, nil
}
