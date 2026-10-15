/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package agent_models

import "beaver/common/models"

// AgentMcpServer 后台配的 MCP 服务。表 agent_mcp_servers。
// 区别于运行时上报的 service/client 级 MCP（agent 进程内缓存，不落库）。
// 本表是配置态，落库 platform，AI 应用引用。
type AgentMcpServer struct {
	models.Model
	McpServerID string `gorm:"column:mcp_server_id;type:varchar(64);uniqueIndex;not null;comment:MCP服务业务ID" json:"mcpServerId"`
	Name        string `gorm:"type:varchar(128);not null;comment:服务名" json:"name"`
	ServerName  string `gorm:"column:server_name;type:varchar(128);not null;index;comment:服务标识 agent进程内缓存按此覆盖" json:"serverName"`
	Transport   string `gorm:"type:varchar(16);not null;default:'stdio';comment:stdio|sse|websocket" json:"transport"`
	Config      string `gorm:"type:json;comment:连接配置 {command/url/env等}" json:"config"`
	Tools       string `gorm:"type:json;comment:能力清单" json:"tools"`
	Status      int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
}

func (AgentMcpServer) TableName() string { return "agent_mcp_servers" }
