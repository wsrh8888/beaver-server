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

package agent_models

import "beaver/common/models"

// AgentApp AI 应用（渠道无关通用配置）。表 agent_apps。
// 对标 Coze Bot / 百炼应用：可被数字人渠道、AI 窗口渠道、未来渠道共同引用。
// 不放任何渠道属性（无 AppKey、无 Runtime、无设备、无工作空间）——渠道属性各归各家。
type AgentApp struct {
	models.Model
	AppID          string `gorm:"column:app_id;type:varchar(64);uniqueIndex;not null;comment:AI应用业务ID" json:"appId"`
	Name           string `gorm:"type:varchar(128);not null;comment:应用名" json:"name"`
	Scene          string `gorm:"type:varchar(64);not null;default:'';index;comment:场景标签 digital_builtin数字伙伴内置/window_defaultAI窗口默认/chat聊天/moment朋友圈" json:"scene"`
	SystemPrompt   string `gorm:"type:text;comment:人设提示词" json:"systemPrompt"`
	LLMModelID     string `gorm:"column:llm_model_id;type:varchar(64);not null;default:'';index;comment:引用 agent_models.model_id" json:"llmModelId"`
	ModelParams    string `gorm:"type:json;comment:模型参数 {temperature,max_turns,...}" json:"modelParams"`
	Tools          string `gorm:"type:json;comment:工具ID列表" json:"tools"`
	McpServerIDs   string `gorm:"column:mcp_server_ids;type:json;comment:后台配的MCP服务ID列表" json:"mcpServerIds"`
	KnowledgeIDs   string `gorm:"column:knowledge_ids;type:json;comment:知识库ID列表" json:"knowledgeIds"`
	OpeningMessage string `gorm:"type:text;comment:开场白" json:"openingMessage"`
	OpeningQuestions string `gorm:"type:json;comment:开场问题" json:"openingQuestions"`
	Variables      string `gorm:"type:json;comment:表单变量" json:"variables"`
	Status         int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
	Sort           int    `gorm:"type:int;not null;default:0;comment:排序" json:"sort"`
	CreatedBy      string `gorm:"column:created_by;type:varchar(32);not null;default:'system';comment:创建来源 portal/backend/system" json:"createdBy"`
}

func (AgentApp) TableName() string { return "agent_apps" }
