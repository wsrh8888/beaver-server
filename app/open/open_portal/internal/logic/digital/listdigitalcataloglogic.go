/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package digital

import (
	"context"
	"errors"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListDigitalCatalogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListDigitalCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDigitalCatalogLogic {
	return &ListDigitalCatalogLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_digital_catalog", ctx)}
}

// ListDigitalCatalog 工具广场与 MCP 目录：平台启用中的工具/服务。
// 数字伙伴只存 ID 引用，名称与描述从这里取，避免在 open 库里复制一份。
func (l *ListDigitalCatalogLogic) ListDigitalCatalog(_ *types.ListDigitalCatalogReq) (resp *types.ListDigitalCatalogRes, err error) {
	toolRes, err := l.svcCtx.PlatformRpc.ListAgentTools(l.ctx, &platform.ListAgentToolsReq{Status: 1})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询工具目录失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("查询工具目录失败")
	}
	mcpRes, err := l.svcCtx.PlatformRpc.ListAgentMcpServers(l.ctx, &platform.ListAgentMcpServersReq{Status: 1})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询 MCP 目录失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("查询 MCP 目录失败")
	}

	tools := make([]types.DigitalToolItem, 0, len(toolRes.Tools))
	for _, item := range toolRes.Tools {
		tools = append(tools, types.DigitalToolItem{
			ToolID:   item.ToolId,
			Name:     item.Name,
			ToolType: item.ToolType,
			Source:   item.Source,
			Status:   int(item.Status),
			Config:   item.Config,
		})
	}
	servers := make([]types.DigitalMcpItem, 0, len(mcpRes.Servers))
	for _, item := range mcpRes.Servers {
		servers = append(servers, types.DigitalMcpItem{
			McpServerID: item.McpServerId,
			Name:        item.Name,
			ServerName:  item.ServerName,
			Transport:   item.Transport,
			Source:      item.Source,
			Status:      int(item.Status),
			Config:      item.Config,
			Tools:       item.Tools,
		})
	}
	return &types.ListDigitalCatalogRes{Tools: tools, McpServers: servers}, nil
}
