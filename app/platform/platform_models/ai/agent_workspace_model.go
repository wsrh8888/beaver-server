/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package ai_models

import "beaver/common/models"

// AgentWorkspace 工作空间注册（AI 助手桌面端特有）。表 agent_workspaces。
// 从 Python agent_models.agent_workspace 迁入。
// 桌面端打开文件夹时上报，agent 跑 LLM 时只读（取当前工作目录），不写。
type AgentWorkspace struct {
	models.Model
	ID        string `gorm:"column:id;type:varchar(64);primaryKey;not null;comment:workspaceId 客户端spaceId" json:"id"`
	UserID    string `gorm:"column:user_id;type:varchar(64);not null;index;comment:所属用户" json:"userId"`
	DeviceID  string `gorm:"column:device_id;type:varchar(64);not null;index;comment:本机设备唯一值" json:"deviceId"`
	Name      string `gorm:"type:varchar(128);not null;default:'';comment:空间名" json:"name"`
	Path      string `gorm:"type:varchar(1024);not null;default:'';comment:本机绝对路径" json:"path"`
	Version   int64  `gorm:"type:bigint;not null;default:0;index;comment:增量同步序列号" json:"version"`
}

func (AgentWorkspace) TableName() string { return "agent_workspaces" }
