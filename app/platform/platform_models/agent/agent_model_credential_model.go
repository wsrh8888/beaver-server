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
	"beaver/common/models"
)

// AgentModelCredential 模型「接入通道」（模型来源）：一条厂商连接 = base URL + API Key。
// 表 agent_model_credentials。对标 stratum 两层货架的上层：多个上架模型共享同一条通道，
// 密钥集中在此、加密落库、只进不出（对外只给 KeyHint），运行时由 GetAgentModel JOIN 后解密下发。
//
// scope=platform：平台官方通道（UserID 为 NULL，全员可用）；
// scope=user：用户个人通道（UserID 必填，仅本人可见），个人自定义模型绑自己的通道。
type AgentModelCredential struct {
	models.Model
	CredentialID    string  `gorm:"column:credential_id;type:varchar(64);uniqueIndex;not null;comment:通道业务ID" json:"credentialId"`
	Name            string  `gorm:"type:varchar(128);not null;comment:通道展示名" json:"name"`
	ModelService    string  `gorm:"column:model_service;type:varchar(64);not null;comment:厂商族 openai|azure|anthropic|deepseek|openai_compatible" json:"modelService"`
	ChannelType     string  `gorm:"column:channel_type;type:varchar(32);not null;default:'openai_compatible';comment:接入类型 official|openai_compatible|azure" json:"channelType"`
	BaseURL         string  `gorm:"column:base_url;type:varchar(500);not null;comment:接口地址" json:"baseUrl"`
	APIKeyEncrypted string  `gorm:"column:api_key_encrypted;type:text;comment:API Key 密文(v1:AES-GCM)" json:"-"`
	KeyHint         string  `gorm:"column:key_hint;type:varchar(32);comment:API Key 脱敏回显" json:"keyHint"`
	Scope           string  `gorm:"type:varchar(16);not null;default:'platform';index;comment:platform官方 user个人" json:"scope"`
	UserID          *string `gorm:"column:user_id;type:varchar(64);index;comment:个人通道所属用户。官方为NULL" json:"userId"`
	Sort            int     `gorm:"type:int;not null;default:0;index;comment:排序" json:"sort"`
	Status          int     `gorm:"type:tinyint;not null;default:1;index;comment:1启用 2停用" json:"status"`
}

func (AgentModelCredential) TableName() string { return "agent_model_credentials" }
