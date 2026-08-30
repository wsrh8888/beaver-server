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

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"beaver/common/models"
)

// 模型档位（对标 OpenAI Luna/Terra/Sol、Claude Haiku/Sonnet/Opus）。
// 表示能力/延迟档，不是价格字段。官方与用户自定义共用。
const (
	ModelTierFast     = "fast"     // 高并发、分类/路由/摘要
	ModelTierBalanced = "balanced" // 默认生产对话与常规工具
	ModelTierStrong   = "strong"   // 深推理、复杂 Agent
)

// AgentModelCapabilities 模型固有能力画像（对标 Azure/OpenAI model card）。
// 官方与用户自定义共用：描述「能做什么」，供任务路由匹配。
type AgentModelCapabilities struct {
	Vision           bool `json:"vision"`           // 多模态视觉
	Tools            bool `json:"tools"`            // 工具/函数调用
	Reasoning        bool `json:"reasoning"`        // 强推理
	StructuredOutput bool `json:"structuredOutput"` // 结构化 JSON 输出
	LongContext      bool `json:"longContext"`      // 长上下文
}

func (c AgentModelCapabilities) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

func (c *AgentModelCapabilities) Scan(value interface{}) error {
	if value == nil {
		*c = AgentModelCapabilities{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return errors.New("type assertion to []byte/string failed")
	}
	if len(b) == 0 {
		*c = AgentModelCapabilities{}
		return nil
	}
	return json.Unmarshal(b, c)
}

// AgentUserModel 用户自定义 LLM 配置（挂在 UserID 下，多端同步）。
// 与官方目录同构：Tier + Capabilities；独有：UserID / Endpoint / ApiKey。
type AgentUserModel struct {
	models.Model
	ModelID      string                 `gorm:"size:64;uniqueIndex;not null;comment:模型业务ID" json:"modelId"`
	UserID       string                 `gorm:"size:64;index;not null;comment:所属用户ID" json:"userId"`
	Name         string                 `gorm:"size:128;not null;comment:模型名称" json:"name"`
	Endpoint     string                 `gorm:"size:500;not null;comment:接口地址" json:"endpoint"`
	ApiKey       string                 `gorm:"size:512;not null;comment:API Key(后续加密存储)" json:"apiKey"`
	Tier         string                 `gorm:"size:16;not null;default:balanced;index;comment:档位fast/balanced/strong" json:"tier"`
	Capabilities AgentModelCapabilities `gorm:"type:longtext;comment:能力画像JSON" json:"capabilities"`
	Status       int8                   `gorm:"type:tinyint;not null;default:1;index;comment:0停用1启用" json:"status"`
}

func (AgentUserModel) TableName() string {
	return "agent_user_models"
}
