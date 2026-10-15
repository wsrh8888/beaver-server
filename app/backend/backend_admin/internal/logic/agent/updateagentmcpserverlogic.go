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

type UpdateAgentMcpServerLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAgentMcpServerLogic {
	return &UpdateAgentMcpServerLogic{logger: beaverlog.New("update_agent_mcp_server", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAgentMcpServerLogic) UpdateAgentMcpServer(req *types.UpdateAgentMcpServerReq) (resp *types.UpdateAgentMcpServerRes, err error) {
	_, err = l.svcCtx.PlatformRpc.UpdateAgentMcpServer(l.ctx, &platform_rpc.UpdateAgentMcpServerReq{
		McpServerId: req.McpServerID,
		Name:        req.Name,
		ServerName:  req.ServerName,
		Transport:   req.Transport,
		Config:      req.Config,
		Tools:       req.Tools,
		Status:      int32(req.Status),
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "更新MCP服务失败", Data: map[string]any{"mcpServerId": req.McpServerID, "err": err.Error()}})
		return nil, err
	}
	return &types.UpdateAgentMcpServerRes{}, nil
}
