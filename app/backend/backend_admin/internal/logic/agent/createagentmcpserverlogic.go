/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package agent

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CreateAgentMcpServerLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAgentMcpServerLogic {
	return &CreateAgentMcpServerLogic{logger: beaverlog.New("create_agent_mcp_server", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateAgentMcpServerLogic) CreateAgentMcpServer(req *types.CreateAgentMcpServerReq) (resp *types.CreateAgentMcpServerRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.CreateAgentMcpServer(l.ctx, &platform_rpc.CreateAgentMcpServerReq{
		Name:       req.Name,
		ServerName: req.ServerName,
		Transport:  req.Transport,
		Config:     req.Config,
		Tools:      req.Tools,
		Status:     int32(req.Status),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "创建MCP服务失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	return &types.CreateAgentMcpServerRes{McpServerID: rpcRes.McpServerId}, nil
}
