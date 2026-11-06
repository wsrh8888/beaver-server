/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-docs/community/license.html
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

// AgentModel 模型目录（上架模型）。表 agent_models。
// 两层货架的下层：本表只存「模型规格 + 绑哪条接入通道」，不再平铺 endpoint/apiKey/provider。
// 厂商、地址、密钥一律由 CredentialID 关联 agent_model_credentials 解析（见 platform_rpc GetAgentModel JOIN）。
// 官方(source=official)与用户自定义(source=custom)同表，用 Source+UserID 区分。
type AgentModel struct {
	models.Model
	ModelID            string      `gorm:"column:model_id;type:varchar(64);uniqueIndex;not null;comment:模型业务ID" json:"modelId"`
	Source             string      `gorm:"type:varchar(16);not null;default:'custom';index;comment:official官方 custom用户自定义" json:"source"`
	UserID             *string     `gorm:"column:user_id;type:varchar(64);index;comment:自定义模型所属用户。官方为NULL" json:"userId"`
	Name               string      `gorm:"type:varchar(128);not null;comment:展示名称" json:"name"`
	ModelName          *string     `gorm:"column:model_name;type:varchar(128);comment:实际调用的模型名" json:"modelName"`
	CredentialID       string      `gorm:"column:credential_id;type:varchar(64);index;comment:关联接入通道 credential_id。provider/endpoint/api_key 由通道解析" json:"credentialId"`
	ContextWindow      int64       `gorm:"column:context_window;type:bigint;not null;default:0;index;comment:上下文窗口tokens" json:"contextWindow"`
	MaxOutputTokens    int         `gorm:"column:max_output_tokens;type:int;not null;default:0;comment:最大输出tokens" json:"maxOutputTokens"`
	InputModalities    string      `gorm:"column:input_modalities;type:json;comment:输入模态[text|image|audio]" json:"inputModalities"`
	OutputModalities   string      `gorm:"column:output_modalities;type:json;comment:输出模态[text|image|audio]" json:"outputModalities"`
	TokenTiers         *TokenTiers `gorm:"column:token_tiers;type:json;comment:各档token上限 {flash,low,medium,high}" json:"tokenTiers"`
	InstructionRole    string      `gorm:"column:instruction_role;type:varchar(16);not null;default:'system';comment:指令角色 system|developer" json:"instructionRole"`
	SupportsThinking   bool        `gorm:"column:supports_thinking;type:tinyint;not null;default:0;comment:是否支持思考" json:"supportsThinking"`
	ThinkingCanDisable bool        `gorm:"column:thinking_can_disable;type:tinyint;not null;default:0;comment:思考是否可关闭" json:"thinkingCanDisable"`
	ThinkingLevels     string      `gorm:"column:thinking_levels;type:json;comment:思考档位[low|medium|high]" json:"thinkingLevels"`
	SupportsTools      bool        `gorm:"column:supports_tools;type:tinyint;not null;default:0;index;comment:是否支持工具调用" json:"supportsTools"`
	SupportsVision     bool        `gorm:"column:supports_vision;type:tinyint;not null;default:0;index;comment:是否支持视觉" json:"supportsVision"`
	Sort               int         `gorm:"type:int;not null;default:0;index;comment:官方模型排序" json:"sort"`
	Status             int         `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
}

func (AgentModel) TableName() string { return "agent_models" }

// TokenTiers AgentModel.TokenTiers 的 JSON 结构：各档 token 上限。
// 实现 driver.Valuer / sql.Scanner，GORM 自动 marshal/unmarshal，列类型用 *TokenTiers 而非裸 string。
// RPC 层仍是 4 个标量字段（对外接口稳定），DB 层合 JSON 避免加档改表。
type TokenTiers struct {
	Flash  int `json:"flash"`
	Low    int `json:"low"`
	Medium int `json:"medium"`
	High   int `json:"high"`
}

// Value 入库：序列化为 JSON 字符串。
func (t TokenTiers) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

// Scan 出库：从 DB 的 JSON 字节反序列化。
func (t *TokenTiers) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	var b []byte
	switch v := val.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return errors.New("TokenTiers.Scan: unsupported type")
	}
	return json.Unmarshal(b, t)
}
