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

package ai_models

import "beaver/common/models"

const (
	SceneAssistant = "assistant" // AI 助手
	SceneChat      = "chat"      // 聊天
	SceneMoment    = "moment"    // 朋友圈
	SceneWorkbench = "workbench" // 工作台
)

// ValidScene 场景写死，每个场景只能有一个智能体。
func ValidScene(scene string) bool {
	switch scene {
	case SceneAssistant, SceneChat, SceneMoment, SceneWorkbench:
		return true
	default:
		return false
	}
}

// AiAssistant 场景智能体。表 ai_assistant。
// scene 唯一：assistant / chat / moment / workbench 各一行。不是 IM 用户，客户端界面单独做。
type AiAssistant struct {
	models.Model
	Scene          string `gorm:"column:scene;type:varchar(32);uniqueIndex;not null;comment:场景 assistant/chat/moment/workbench" json:"scene"`
	Name           string `gorm:"type:varchar(128);not null;default:'';comment:助手名称" json:"name"`
	SystemPrompt   string `gorm:"column:system_prompt;type:text;comment:人设提示词" json:"systemPrompt"`
	LLMModelID     string `gorm:"column:llm_model_id;type:varchar(64);not null;default:'';comment:引用 agent_models.model_id" json:"llmModelId"`
	ModelParams    string `gorm:"column:model_params;type:json;comment:模型参数" json:"modelParams"`
	Tools          string `gorm:"type:json;comment:工具ID列表" json:"tools"`
	McpServerIDs   string `gorm:"column:mcp_server_ids;type:json;comment:MCP服务ID列表" json:"mcpServerIds"`
	OpeningMessage string `gorm:"column:opening_message;type:text;comment:开场白" json:"openingMessage"`
	Status         int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
}

func (AiAssistant) TableName() string { return "ai_assistant" }
