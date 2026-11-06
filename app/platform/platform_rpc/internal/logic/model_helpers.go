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

package logic

import (
	agent_models "beaver/app/platform/platform_models/agent"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/utils/secret"

	"gorm.io/gorm"
)

// toAgentModelItem 把「模型行 + 其绑定的接入通道」拼成对外 AgentModelItem。
// revealKey=true（运行时 GetAgentModel）：解密通道 Key 填 api_key，供 agent 直调厂商；
// revealKey=false（后台/客户端 ListAgentModels）：绝不透出真实 Key，只给通道名/类型，Key 恒空。
func toAgentModelItem(
	row *agent_models.AgentModel,
	cred *agent_models.AgentModelCredential,
	cipher *secret.Cipher,
	revealKey bool,
) *platform_rpc.AgentModelItem {
	if row == nil {
		return nil
	}
	var tiers agent_models.TokenTiers
	if row.TokenTiers != nil {
		tiers = *row.TokenTiers
	}
	item := &platform_rpc.AgentModelItem{
		Id:                 uint64(row.Id),
		ModelId:            row.ModelID,
		Source:             row.Source,
		Name:               row.Name,
		ModelName:          strVal(row.ModelName),
		ContextWindow:      row.ContextWindow,
		MaxOutputTokens:    int32(row.MaxOutputTokens),
		MaxTokensFlash:     int32(tiers.Flash),
		MaxTokensLow:       int32(tiers.Low),
		MaxTokensMedium:    int32(tiers.Medium),
		MaxTokensHigh:      int32(tiers.High),
		InstructionRole:    row.InstructionRole,
		SupportsThinking:   row.SupportsThinking,
		ThinkingCanDisable: row.ThinkingCanDisable,
		SupportsTools:      row.SupportsTools,
		SupportsVision:     row.SupportsVision,
		Sort:               int32(row.Sort),
		Status:             int32(row.Status),
		InputModalities:    row.InputModalities,
		OutputModalities:   row.OutputModalities,
		ThinkingLevels:     row.ThinkingLevels,
		CredentialId:       row.CredentialID,
	}
	if row.UserID != nil {
		item.UserId = *row.UserID
	}
	if cred != nil {
		// provider/endpoint 由通道解析而来；channel_name/type 供后台回显「这个模型走的哪条通道」。
		item.Provider = cred.ModelService
		item.Endpoint = cred.BaseURL
		item.ChannelName = cred.Name
		item.ChannelType = cred.ChannelType
		if revealKey && cipher != nil && cred.APIKeyEncrypted != "" {
			if plain, err := cipher.Decrypt(cred.APIKeyEncrypted); err == nil {
				item.ApiKey = plain
			}
		}
	}
	return item
}

// fetchCredential 按业务 ID 取单条通道；取不到返回 nil（调用方按「未绑定通道」处理）。
func fetchCredential(db *gorm.DB, credentialID string) *agent_models.AgentModelCredential {
	if credentialID == "" {
		return nil
	}
	var cred agent_models.AgentModelCredential
	if err := db.Where("credential_id = ?", credentialID).First(&cred).Error; err != nil {
		return nil
	}
	return &cred
}

// fetchCredentialMap 批量取通道，避免列表接口逐行 N+1 查询。空 ID 自动跳过。
func fetchCredentialMap(db *gorm.DB, ids []string) map[string]*agent_models.AgentModelCredential {
	out := make(map[string]*agent_models.AgentModelCredential)
	uniq := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	if len(uniq) == 0 {
		return out
	}
	var creds []agent_models.AgentModelCredential
	db.Where("credential_id IN ?", uniq).Find(&creds)
	for i := range creds {
		c := creds[i]
		out[c.CredentialID] = &c
	}
	return out
}

// toCredentialItem 拼接入通道对外项。只透 key_hint + key_configured，绝不透密文或明文。
func toCredentialItem(c *agent_models.AgentModelCredential) *platform_rpc.AgentCredentialItem {
	if c == nil {
		return nil
	}
	item := &platform_rpc.AgentCredentialItem{
		Id:            uint64(c.Id),
		CredentialId:  c.CredentialID,
		Name:          c.Name,
		ModelService:  c.ModelService,
		ChannelType:   c.ChannelType,
		BaseUrl:       c.BaseURL,
		KeyHint:       c.KeyHint,
		KeyConfigured: c.APIKeyEncrypted != "",
		Scope:         c.Scope,
		Sort:          int32(c.Sort),
		Status:        int32(c.Status),
	}
	if c.UserID != nil {
		item.UserId = *c.UserID
	}
	return item
}

// validCredentialBinding 校验模型与接入通道的绑定是否合法：
//   - 通道必须存在且启用；
//   - source=official → 必须绑平台(scope=platform)通道；
//   - source=custom   → 必须绑本人的(scope=user 且 user_id 相符)通道。
func validCredentialBinding(db *gorm.DB, credentialID, source, userID string) bool {
	if credentialID == "" {
		return false
	}
	cred := fetchCredential(db, credentialID)
	if cred == nil || cred.Status != 1 {
		return false
	}
	if source == "custom" {
		return cred.Scope == "user" && userID != "" && cred.UserID != nil && *cred.UserID == userID
	}
	return cred.Scope == "platform"
}

func strVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
