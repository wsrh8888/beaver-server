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

type DeleteDigitalMcpLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteDigitalMcpLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDigitalMcpLogic {
	return &DeleteDigitalMcpLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("delete_digital_mcp", ctx)}
}

// DeleteDigitalMcp 删除 MCP 服务，并从本数字伙伴的启用列表里摘掉，
// 否则 config 里留着一个已删的 ID，运行时挂载时会报错。
func (l *DeleteDigitalMcpLogic) DeleteDigitalMcp(req *types.DeleteDigitalMcpReq) (resp *types.DeleteDigitalMcpRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	partner, err := loadOwnedDigital(l.svcCtx.DB, req.UserID, req.DigitalID)
	if err != nil {
		return nil, err
	}
	if req.McpServerID == "" {
		return nil, errors.New("mcpServerId 不能为空")
	}
	if _, err := l.svcCtx.PlatformRpc.DeleteAgentMcpServer(l.ctx, &platform.DeleteAgentMcpServerReq{
		McpServerId: req.McpServerID,
	}); err != nil {
		l.logger.Error(model.LogMsg{Text: "删除 MCP 服务失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, errors.New("删除 MCP 服务失败")
	}

	ids := loadDigitalMcpServerIDs(l.svcCtx.DB, partner.DigitalID)
	remain := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != req.McpServerID {
			remain = append(remain, id)
		}
	}
	if err := writeDigitalMcpServerIDs(l.svcCtx.DB, partner.DigitalID, remain); err != nil {
		l.logger.Error(model.LogMsg{Text: "解绑已删除的 MCP 服务失败", Data: map[string]interface{}{"err": err.Error()}})
	}
	return &types.DeleteDigitalMcpRes{}, nil
}
