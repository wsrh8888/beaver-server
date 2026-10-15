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

type GetAgentMcpServerLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentMcpServerLogic {
	return &GetAgentMcpServerLogic{logger: beaverlog.New("get_agent_mcp_server", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentMcpServerLogic) GetAgentMcpServer(req *types.GetAgentMcpServerReq) (resp *types.GetAgentMcpServerRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.GetAgentMcpServer(l.ctx, &platform_rpc.GetAgentMcpServerReq{
		McpServerId: req.McpServerID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取MCP服务失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	out := types.GetAgentMcpServerRes{}
	if rpcRes.Server != nil {
		out.Server = toAgentMcpServerInfo(rpcRes.Server)
	}
	return &out, nil
}
