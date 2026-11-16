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
// 手写 JSON DTO：字段名与 json tag 与对端 beaver-open 的 .api 定义（go-zero httpx.Parse）逐字段一致。
// 对端 .api 一律用 lowerCamelCase，故本包 tag 全部驼峰——写成下划线会导致对端 httpx.Parse
// 认为必填字段未设置（field "xxx" is not set），或响应解析出零值。
package open

type BotItem struct {
	Id               uint32 `json:"id,omitempty"`
	BotId            string `json:"botId,omitempty"`
	OwnerId          string `json:"ownerId,omitempty"`
	Visibility       string `json:"visibility,omitempty"`
	Status           int32  `json:"status,omitempty"`
	SignatureEnabled int32  `json:"signatureEnabled,omitempty"`
	CreatedAt        int64  `json:"createdAt,omitempty"`
}

type BotSecurity struct {
	KeywordsEnabled    bool     `json:"keywordsEnabled,omitempty"`
	Keywords           []string `json:"keywords,omitempty"`
	IpWhitelistEnabled bool     `json:"ipWhitelistEnabled,omitempty"`
	IpWhitelist        []string `json:"ipWhitelist,omitempty"`
}

type CreateDigitalReq struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	Runtime      string `json:"runtime,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
	Source       string `json:"source,omitempty"`
	OwnerId      string `json:"ownerId,omitempty"`
	SystemPrompt string `json:"systemPrompt,omitempty"`
	LlmModelId   string `json:"llmModelId,omitempty"`
	Tools        string `json:"tools,omitempty"`
	McpServerIds string `json:"mcpServerIds,omitempty"`
}

type CreateDigitalRes struct {
	DigitalId     string `json:"digitalId,omitempty"`
	DigitalKey    string `json:"digitalKey,omitempty"`
	DigitalSecret string `json:"digitalSecret,omitempty"`
}

type DigitalItem struct {
	DigitalId    string `json:"digitalId,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	OwnerId      string `json:"ownerId,omitempty"`
	Runtime      string `json:"runtime,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
	Status       int32  `json:"status,omitempty"`
	Source       string `json:"source,omitempty"`
	CreatedAt    int64  `json:"createdAt,omitempty"`
	SystemPrompt string `json:"systemPrompt,omitempty"`
	LlmModelId   string `json:"llmModelId,omitempty"`
	Tools        string `json:"tools,omitempty"`
	McpServerIds string `json:"mcpServerIds,omitempty"`
}

type ExchangeTokenReq struct {
	OauthId string `json:"oauthId,omitempty"`
	Code    string `json:"code,omitempty"`
}

type ExchangeTokenRes struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	ExpiresIn    int64  `json:"expiresIn,omitempty"`
	TokenType    string `json:"tokenType,omitempty"`
	UserId       string `json:"userId,omitempty"`
}

type GetDigitalByUserIDReq struct {
	DigitalUserId string `json:"digitalUserId,omitempty"`
}

type GetDigitalByUserIDRes struct {
	Found         bool   `json:"found,omitempty"`
	DigitalKey    string `json:"digitalKey,omitempty"`
	DigitalUserId string `json:"digitalUserId,omitempty"`
	Runtime       string `json:"runtime,omitempty"`
}

type GetDigitalConfigReq struct {
	DigitalId string `json:"digitalId"`
}

type GetDigitalConfigRes struct {
	DigitalId string `json:"digitalId"`
	Runtime   string `json:"runtime"`
	// AgentId 数字伙伴绑定的 Agent；runtime=builtin/workflow 时非空，external 时为空。
	// 编排配置（人设/模型/工具/知识）只在 agents 表，open 侧不再下发编排 JSON。
	AgentId        string `json:"agentId"`
	ExternalConfig string `json:"externalConfig"`
}

type ListBotsReq struct {
	Page        int32  `json:"page,omitempty"`
	PageSize    int32  `json:"pageSize,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	Status      int32  `json:"status,omitempty"`
	OwnerUserId string `json:"ownerUserId,omitempty"`
}

type ListBotsRes struct {
	Total int64      `json:"total,omitempty"`
	List  []*BotItem `json:"list,omitempty"`
}

type ListDigitalsReq struct {
	Source   string `json:"source,omitempty"`
	Status   int32  `json:"status,omitempty"`
	Page     int32  `json:"page,omitempty"`
	PageSize int32  `json:"pageSize,omitempty"`
}

type ListDigitalsRes struct {
	Total int64          `json:"total,omitempty"`
	List  []*DigitalItem `json:"list,omitempty"`
}

type ListOpenOAuthReq struct {
	Page        int32  `json:"page,omitempty"`
	PageSize    int32  `json:"pageSize,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	OwnerUserId string `json:"ownerUserId,omitempty"`
	Status      int32  `json:"status,omitempty"`
	OauthId     string `json:"oauthId,omitempty"`
}

type ListOpenOAuthRes struct {
	Total int64            `json:"total,omitempty"`
	List  []*OpenOAuthItem `json:"list,omitempty"`
}

type OpenOAuthItem struct {
	OauthId     string `json:"oauthId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	OwnerUserId string `json:"ownerUserId,omitempty"`
	Status      int32  `json:"status,omitempty"`
	CreatedAt   int64  `json:"createdAt,omitempty"`
}

type UpdateBotReq struct {
	BotId            string       `json:"botId,omitempty"`
	Visibility       string       `json:"visibility,omitempty"`
	SignatureEnabled int32        `json:"signatureEnabled,omitempty"`
	Security         *BotSecurity `json:"security,omitempty"`
	Status           int32        `json:"status,omitempty"`
}

type UpdateBotRes struct{}

type UpdateDigitalReq struct {
	DigitalId    string `json:"digitalId,omitempty"`
	Name         string `json:"name,omitempty"`
	Visibility   string `json:"visibility,omitempty"`
	Status       int32  `json:"status,omitempty"`
	SystemPrompt string `json:"systemPrompt,omitempty"`
	LlmModelId   string `json:"llmModelId,omitempty"`
	Tools        string `json:"tools,omitempty"`
	McpServerIds string `json:"mcpServerIds,omitempty"`
}

type UpdateDigitalRes struct{}

type UpdateOpenOAuthReq struct {
	OauthIds   []string `json:"oauthIds,omitempty"`
	Status     int32    `json:"status,omitempty"`
	OperatorId string   `json:"operatorId,omitempty"`
}

type UpdateOpenOAuthRes struct {
	AffectedCount int64 `json:"affectedCount,omitempty"`
}
