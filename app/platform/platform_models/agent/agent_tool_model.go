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

// AgentTool 工具/能力目录。表 agent_tools。
// 后台配、AI 应用引用，跨渠道通用。
type AgentTool struct {
	models.Model
	ToolID    string `gorm:"column:tool_id;type:varchar(64);uniqueIndex;not null;comment:工具业务ID" json:"toolId"`
	Name      string `gorm:"type:varchar(128);not null;comment:工具名" json:"name"`
	ToolType  string `gorm:"column:tool_type;type:varchar(32);not null;default:'builtin';comment:builtin|mcp|api" json:"toolType"`
	Config    string `gorm:"type:json;comment:工具配置" json:"config"`
	Status    int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
	Sort      int    `gorm:"type:int;not null;default:0;comment:排序" json:"sort"`
}

func (AgentTool) TableName() string { return "agent_tools" }
