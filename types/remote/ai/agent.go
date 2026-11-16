/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

// Package ai AI 服务（beaver-ai）S2S 出入参类型，与对端 ai_rpc proto JSON 逐字段对齐。
//
// 放这里而非 remote/ai 内部，便于调用方只依赖类型、不依赖客户端实现。
package ai

// GetAgentReq 取 Agent 编排配置请求（agent_id 与 scene 二选一，agent_id 优先）。
type GetAgentReq struct {
	AgentID string `json:"agent_id"`
	Scene   string `json:"scene"`
}

// AgentItem Agent 编排配置项（agents 表一行，渠道无关）。
type AgentItem struct {
	AgentID      string `json:"agent_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Source       string `json:"source"`
	Runtime      string `json:"runtime"`
	SystemPrompt string `json:"system_prompt"`
	LlmModelID   string `json:"llm_model_id"`
}

// GetAgentRes 取 Agent 编排配置响应。
type GetAgentRes struct {
	Agent *AgentItem `json:"agent"`
}
