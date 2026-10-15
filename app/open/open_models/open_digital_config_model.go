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

package open_models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// OpenDigitalConfig 数字伙伴编排配置（每个数字伙伴一行）。
//
// 数据归属在 open 模块（用户在 portal 编排数字伙伴时填写），agent 模块运行时通过 RPC 来取。
// 按 Runtime 分列独立 JSON 字段：builtin/external/workflow 各自一个 JSON 字符串，互不污染，
// 模式差异不挤在一个 struct 里；新增模式只加一列 JSON，不改已有列。
type OpenDigitalConfig struct {
	gorm.Model
	DigitalID string `gorm:"column:digital_id;type:varchar(64);uniqueIndex;not null;comment:数字伙伴ID（OpenDigital.DigitalID）"`

	AppID         string            `gorm:"column:app_id;type:varchar(64);not null;default:'';index;comment:引用 agent_apps.app_id（builtin 模式）"`
	OverrideConfig string           `gorm:"column:override_config;type:json;comment:覆盖 AI 应用部分字段的 JSON（可空）"`
	BuiltinConfig   *BuiltinConfig   `gorm:"column:builtin_config;type:json;comment:提示词模式配置 JSON（Runtime=builtin）"`
	ExternalConfig  *ExternalConfig  `gorm:"column:external_config;type:json;comment:长链接模式配置 JSON（Runtime=external）"`
	WorkflowConfig  *WorkflowConfig  `gorm:"column:workflow_config;type:json;comment:工作流模式配置 JSON（Runtime=workflow）"`
}

func (OpenDigitalConfig) TableName() string { return "open_digital_config" }

// BuiltinConfig 提示词模式配置（Runtime=builtin）。
type BuiltinConfig struct {
	SystemPrompt    string            `json:"systemPrompt"`              // 人设提示词
	LLMModelID       string            `json:"llmModelId"`               // 模型ID，引用 agent_models.model_id
	ModelParams      map[string]any    `json:"modelParams,omitempty"`     // 模型参数 temperature、max_turns 等
	OpeningMessage   string            `json:"openingMessage,omitempty"`  // 对话开场白，空=不主动开口
	OpeningQuestions []string          `json:"openingQuestions,omitempty"` // 开场问题
	Variables       []ConfigVariable  `json:"variables,omitempty"`        // 表单变量（用户对话前填写，替换 {{variable}}）
	Tools           []string           `json:"tools,omitempty"`            // 启用的工具ID列表
	KnowledgeIDs    []string           `json:"knowledgeIds,omitempty"`      // 绑定的知识库ID列表，跨服务只存ID
}

// ExternalConfig 长链接模式配置（Runtime=external）。
type ExternalConfig struct {
	CallbackURL string `json:"callbackUrl,omitempty"` // 外部 agent 回调地址（按需）
}

// WorkflowConfig 工作流模式配置（Runtime=workflow）。
type WorkflowConfig struct {
	Nodes []WorkflowNode `json:"nodes"` // 流程节点（start/llm/condition/tool/knowledge/end）
	Edges []WorkflowEdge `json:"edges"` // 节点连线（含条件分支）
}

// WorkflowNode 工作流节点。
//
// NodeType 决定 NodeConfig 的结构，各节点配置整体落 NodeConfig JSON，不拆结构化字段——
// 新增节点类型只加 NodeType 常量 + NodeConfig 字段约定，不改表结构。
type WorkflowNode struct {
	NodeID     string         `json:"nodeId"`               // 节点ID（流程内唯一）
	NodeType   string         `json:"nodeType"`             // start|llm|condition|tool|knowledge|end
	NodeConfig map[string]any `json:"nodeConfig,omitempty"` // 节点配置（随 NodeType 不同）
}

// WorkflowEdge 工作流连线。
type WorkflowEdge struct {
	EdgeID    string `json:"edgeId"`             // 连线ID
	FromNode  string `json:"fromNode"`           // 起点节点ID
	ToNode    string `json:"toNode"`             // 终点节点ID
	// 条件分支：condition 节点出边按 Condition 匹配决定走哪条；普通节点出边 Condition 为空。
	Condition string `json:"condition,omitempty"` // 条件表达式（condition 节点分支用）
}

// ConfigVariable 表单变量（用户在对话前填写，替换提示词中的 {{variable}}）。
type ConfigVariable struct {
	Name     string   `json:"name"`               // 变量名（提示词中 {{name}} 引用）
	Label    string   `json:"label"`              // 展示标签
	Type     string   `json:"type"`               // text|number|select
	Default  string   `json:"default,omitempty"`
	Required bool     `json:"required,omitempty"`
	Options  []string `json:"options,omitempty"` // type=select 时的选项
}

// —— Scan/Value：参考 ctype.Msg，三个配置各自独立序列化 ——

func (c *BuiltinConfig) Scan(val any) error {
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}
func (c BuiltinConfig) Value() (driver.Value, error) { return json.Marshal(c) }

func (c *ExternalConfig) Scan(val any) error {
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}
func (c ExternalConfig) Value() (driver.Value, error) { return json.Marshal(c) }

func (c *WorkflowConfig) Scan(val any) error {
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}
func (c WorkflowConfig) Value() (driver.Value, error) { return json.Marshal(c) }
