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

// beaver-open S2S 出入参类型
//
// 手写 JSON DTO：字段名与 json tag 与对端 proto 的 JSON 编码逐字段一致，
// 不依赖 protobuf 运行时（无 protoimpl / protoreflect / state）。
package open

type BotItem struct {
	Id               uint32 `json:"id,omitempty"`
	BotId            string `json:"bot_id,omitempty"`
	OwnerId          string `json:"owner_id,omitempty"`
	Visibility       string `json:"visibility,omitempty"`
	Status           int32  `json:"status,omitempty"`
	SignatureEnabled int32  `json:"signature_enabled,omitempty"`
	CreatedAt        int64  `json:"created_at,omitempty"`
}

type BotSecurity struct {
	KeywordsEnabled    bool     `json:"keywords_enabled,omitempty"`
	Keywords           []string `json:"keywords,omitempty"`
	IpWhitelistEnabled bool     `json:"ip_whitelist_enabled,omitempty"`
	IpWhitelist        []string `json:"ip_whitelist,omitempty"`
}

type CreateDigitalReq struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Runtime      string `json:"runtime,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	Source       string `json:"source,omitempty"`
	OwnerId      string `json:"owner_id,omitempty"`
	SystemPrompt string `json:"system_prompt,omitempty"`
	LlmModelId   string `json:"llm_model_id,omitempty"`
	Tools        string `json:"tools,omitempty"`
	McpServerIds string `json:"mcp_server_ids,omitempty"`
}

type CreateDigitalRes struct {
	DigitalId string `json:"digital_id,omitempty"`
	AppKey    string `json:"app_key,omitempty"`
	AppSecret string `json:"app_secret,omitempty"`
}

type DigitalItem struct {
	DigitalId    string `json:"digital_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	OwnerId      string `json:"owner_id,omitempty"`
	Runtime      string `json:"runtime,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
	Status       int32  `json:"status,omitempty"`
	Source       string `json:"source,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	SystemPrompt string `json:"system_prompt,omitempty"`
	LlmModelId   string `json:"llm_model_id,omitempty"`
	Tools        string `json:"tools,omitempty"`
	McpServerIds string `json:"mcp_server_ids,omitempty"`
}

type ExchangeTokenReq struct {
	OauthId string `json:"oauth_id,omitempty"`
	Code    string `json:"code,omitempty"`
}

type ExchangeTokenRes struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	UserId       string `json:"user_id,omitempty"`
}

type GetDigitalByUserIDReq struct {
	DigitalUserId string `json:"digital_user_id,omitempty"`
}

type GetDigitalByUserIDRes struct {
	Found         bool   `json:"found,omitempty"`
	AppKey        string `json:"app_key,omitempty"`
	DigitalUserId string `json:"digital_user_id,omitempty"`
	Runtime       string `json:"runtime,omitempty"`
}

type GetDigitalConfigReq struct {
	DigitalId string `json:"digital_id,omitempty"`
}

type GetDigitalConfigRes struct {
	DigitalId      string `json:"digital_id,omitempty"`
	Runtime        string `json:"runtime,omitempty"`
	BuiltinConfig  string `json:"builtin_config,omitempty"`
	ExternalConfig string `json:"external_config,omitempty"`
	WorkflowConfig string `json:"workflow_config,omitempty"`
	AppId          string `json:"app_id,omitempty"`
}

type ListBotsReq struct {
	Page        int32  `json:"page,omitempty"`
	PageSize    int32  `json:"page_size,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	Status      int32  `json:"status,omitempty"`
	OwnerUserId string `json:"owner_user_id,omitempty"`
}

type ListBotsRes struct {
	Total int64      `json:"total,omitempty"`
	List  []*BotItem `json:"list,omitempty"`
}

type ListDigitalsReq struct {
	Source   string `json:"source,omitempty"`
	Status   int32  `json:"status,omitempty"`
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"page_size,omitempty"`
}

type ListDigitalsRes struct {
	Total int64          `json:"total,omitempty"`
	List  []*DigitalItem `json:"list,omitempty"`
}

type ListOpenOAuthReq struct {
	Page        int32  `json:"page,omitempty"`
	PageSize    int32  `json:"page_size,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	OwnerUserId string `json:"owner_user_id,omitempty"`
	Status      int32  `json:"status,omitempty"`
	OauthId     string `json:"oauth_id,omitempty"`
}

type ListOpenOAuthRes struct {
	Total int64            `json:"total,omitempty"`
	List  []*OpenOAuthItem `json:"list,omitempty"`
}

type OpenOAuthItem struct {
	OauthId     string `json:"oauth_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	OwnerUserId string `json:"owner_user_id,omitempty"`
	Status      int32  `json:"status,omitempty"`
	CreatedAt   int64  `json:"created_at,omitempty"`
}

type UpdateBotReq struct {
	BotId            string       `json:"bot_id,omitempty"`
	Visibility       string       `json:"visibility,omitempty"`
	SignatureEnabled int32        `json:"signature_enabled,omitempty"`
	Security         *BotSecurity `json:"security,omitempty"`
	Status           int32        `json:"status,omitempty"`
}

type UpdateBotRes struct{}

type UpdateDigitalReq struct {
	DigitalId    string `json:"digital_id,omitempty"`
	Name         string `json:"name,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
	Status       int32  `json:"status,omitempty"`
	AppId        string `json:"app_id,omitempty"`
	SystemPrompt string `json:"system_prompt,omitempty"`
	LlmModelId   string `json:"llm_model_id,omitempty"`
	Tools        string `json:"tools,omitempty"`
	McpServerIds string `json:"mcp_server_ids,omitempty"`
}

type UpdateDigitalRes struct{}

type UpdateOpenOAuthReq struct {
	OauthIds   []string `json:"oauth_ids,omitempty"`
	Status     int32    `json:"status,omitempty"`
	OperatorId string   `json:"operator_id,omitempty"`
}

type UpdateOpenOAuthRes struct {
	AffectedCount int64 `json:"affected_count,omitempty"`
}
