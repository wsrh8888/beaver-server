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

// GetAgentAppReq 取 AI 应用配置请求。
type GetAgentAppReq struct {
	AppId string `json:"app_id"`
}

// AgentAppItem AI 应用配置项。
type AgentAppItem struct {
	AppId        string `json:"app_id"`
	Name         string `json:"name"`
	Scene        string `json:"scene"`
	SystemPrompt string `json:"system_prompt"`
}

// GetAgentAppRes 取 AI 应用配置响应。
type GetAgentAppRes struct {
	App *AgentAppItem `json:"app"`
}
