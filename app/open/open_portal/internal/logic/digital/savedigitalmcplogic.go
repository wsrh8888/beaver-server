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
	"strings"

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/platform/platform_rpc/platform"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type SaveDigitalMcpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewSaveDigitalMcpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveDigitalMcpLogic {
	return &SaveDigitalMcpLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("save_digital_mcp", ctx)}
}

// SaveDigitalMcp 新建或更新 MCP 服务（平台级资源，数字伙伴只存 ID 引用）。
// 建完不自动绑定：MCP 连接配置填错会直接影响运行时，让用户先看清再启用。
func (l *SaveDigitalMcpLogic) SaveDigitalMcp(req *types.SaveDigitalMcpReq) (resp *types.SaveDigitalMcpRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if _, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("服务名称不能为空")
	}
	transport := strings.TrimSpace(req.Transport)
	if transport == "" {
		transport = "websocket"
	}
	serverName := strings.TrimSpace(req.ServerName)

	if req.McpServerID == "" {
		created, err := l.svcCtx.PlatformRpc.CreateAgentMcpServer(l.ctx, &platform.CreateAgentMcpServerReq{
			Name:       name,
			ServerName: serverName,
			Transport:  transport,
			Config:     req.Config,
			Tools:      req.Tools,
			Status:     1,
			Source:     "custom",
			UserId:     req.UserID,
		})
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "新建 MCP 服务失败", Data: map[string]interface{}{"err": err.Error()}})
			return nil, errors.New("新建 MCP 服务失败")
		}
		return &types.SaveDigitalMcpRes{McpServerID: created.McpServerId}, nil
	}

	if _, err := l.svcCtx.PlatformRpc.UpdateAgentMcpServer(l.ctx, &platform.UpdateAgentMcpServerReq{
		McpServerId: req.McpServerID,
		Name:        name,
		ServerName:  serverName,
		Transport:   transport,
		Config:      req.Config,
		Tools:       req.Tools,
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "更新 MCP 服务失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("更新 MCP 服务失败")
	}
	return &types.SaveDigitalMcpRes{McpServerID: req.McpServerID}, nil
}
