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
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

// toAgentAppInfo platform_rpc.AgentAppItem → 后台 types.AgentAppInfo
func toAgentAppInfo(a *platform_rpc.AgentAppItem) types.AgentAppInfo {
	return types.AgentAppInfo{
		AppID:            a.AppId,
		Name:             a.Name,
		Scene:            a.Scene,
		SystemPrompt:     a.SystemPrompt,
		LLMModelID:       a.LlmModelId,
		ModelParams:      a.ModelParams,
		Tools:            a.Tools,
		McpServerIDs:     a.McpServerIds,
		KnowledgeIDs:     a.KnowledgeIds,
		OpeningMessage:   a.OpeningMessage,
		OpeningQuestions: a.OpeningQuestions,
		Variables:        a.Variables,
		Status:           int(a.Status),
		Sort:             int(a.Sort),
		CreatedBy:        a.CreatedBy,
	}
}

// toAgentToolInfo platform_rpc.AgentToolItem → 后台 types.AgentToolInfo
func toAgentToolInfo(t *platform_rpc.AgentToolItem) types.AgentToolInfo {
	return types.AgentToolInfo{
		ToolID:   t.ToolId,
		Name:     t.Name,
		ToolType: t.ToolType,
		Config:   t.Config,
		Status:   int(t.Status),
		Sort:     int(t.Sort),
	}
}

// toAgentMcpServerInfo platform_rpc.AgentMcpServerItem → 后台 types.AgentMcpServerInfo
func toAgentMcpServerInfo(s *platform_rpc.AgentMcpServerItem) types.AgentMcpServerInfo {
	return types.AgentMcpServerInfo{
		McpServerID: s.McpServerId,
		Name:        s.Name,
		ServerName:  s.ServerName,
		Transport:   s.Transport,
		Config:      s.Config,
		Tools:       s.Tools,
		Status:      int(s.Status),
	}
}
