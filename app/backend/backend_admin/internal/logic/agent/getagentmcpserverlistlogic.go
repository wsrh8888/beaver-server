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

type GetAgentMcpServerListLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAgentMcpServerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAgentMcpServerListLogic {
	return &GetAgentMcpServerListLogic{logger: beaverlog.New("get_agent_mcp_server_list", ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAgentMcpServerListLogic) GetAgentMcpServerList(req *types.GetAgentMcpServerListReq) (resp *types.GetAgentMcpServerListRes, err error) {
	rpcRes, err := l.svcCtx.PlatformRpc.ListAgentMcpServers(l.ctx, &platform_rpc.ListAgentMcpServersReq{
		Status: int32(req.Status),
		Source: req.Source,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取MCP服务列表失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}
	list := make([]types.AgentMcpServerInfo, 0, len(rpcRes.Servers))
	for _, s := range rpcRes.Servers {
		list = append(list, toAgentMcpServerInfo(s))
	}
	return &types.GetAgentMcpServerListRes{List: list}, nil
}
