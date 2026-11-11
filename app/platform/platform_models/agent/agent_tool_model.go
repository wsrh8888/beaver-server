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

// AgentTool 技能目录。表 agent_tools。
//
// 一个技能 = 一份 SKILL.md 文件。上传后按 frontmatter 解析出名称、描述、版本，
// 正文原样存 skill_md，运行时读这份正文决定「什么时候用、怎么用」。
// 后台配、AI 应用按 tool_id 引用，跨渠道通用。
type AgentTool struct {
	models.Model
	ToolID      string `gorm:"column:tool_id;type:varchar(64);uniqueIndex;not null;comment:技能业务ID" json:"toolId"`
	Name        string `gorm:"type:varchar(128);not null;comment:技能名，取自 SKILL.md frontmatter" json:"name"`
	Description string `gorm:"type:varchar(512);not null;default:'';comment:一句话描述，运行时用来判断是否调用" json:"description"`
	Version     string `gorm:"type:varchar(32);not null;default:'';comment:技能版本" json:"version"`
	ToolType    string `gorm:"column:tool_type;type:varchar(32);not null;default:'skill';comment:skill，历史 builtin/mcp/api 保留兼容" json:"toolType"`
	SkillMD     string `gorm:"column:skill_md;type:text;comment:SKILL.md 全文" json:"skillMd"`
	FileName    string `gorm:"column:file_name;type:varchar(255);not null;default:'';comment:上传时的原始文件名" json:"fileName"`
	Config      string `gorm:"type:text;comment:frontmatter 其余字段（allowed-tools 等）" json:"config"`
	Status      int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
	Sort        int    `gorm:"type:int;not null;default:0;comment:排序" json:"sort"`
	Source      string `gorm:"type:varchar(16);not null;default:'official';index;comment:official官方 custom用户" json:"source"`
	UserID      string `gorm:"column:user_id;type:varchar(64);not null;default:'';index;comment:用户自建时的用户ID，官方为空" json:"userId"`
}

func (AgentTool) TableName() string { return "agent_tools" }
