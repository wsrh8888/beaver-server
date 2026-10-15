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

type DeleteAgentMcpServerLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAgentMcpServerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAgentMcpServerLogic {
	return &DeleteAgentMcpServerLogic{logger: beaverlog.New("delete_agent_mcp_server", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteAgentMcpServerLogic) DeleteAgentMcpServer(req *types.DeleteAgentMcpServerReq) (resp *types.DeleteAgentMcpServerRes, err error) {
	_, err = l.svcCtx.PlatformRpc.DeleteAgentMcpServer(l.ctx, &platform_rpc.DeleteAgentMcpServerReq{
		McpServerId: req.McpServerID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "删除MCP服务失败", Data: map[string]any{"mcpServerId": req.McpServerID, "err": err.Error()}})
		return nil, err
	}
	return &types.DeleteAgentMcpServerRes{}, nil
}
