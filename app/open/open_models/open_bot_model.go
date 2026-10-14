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

const (
	BotVisibilityPrivate = "private" // 仅管理员可搜到（Owner + role=admin）
	BotVisibilityPartial = "partial" // 指定人可搜到（role=admin|viewer）
	BotVisibilityPublic  = "public"  // 全 IM 用户可搜到
)

const (
	BotStatusEnabled  = 1
	BotStatusDisabled = 0
)

// OpenBotSecurity 群通知机器人入站校验配置（JSON；不含签名密钥）。
type OpenBotSecurity struct {
	KeywordsEnabled    bool     `json:"keywordsEnabled"`
	Keywords           []string `json:"keywords"`
	IPWhitelistEnabled bool     `json:"ipWhitelistEnabled"`
	IPWhitelist        []string `json:"ipWhitelist"`
}

func (s OpenBotSecurity) Value() (driver.Value, error) {
	if !s.KeywordsEnabled && !s.IPWhitelistEnabled {
		return "{}", nil
	}
	return json.Marshal(s)
}

func (s *OpenBotSecurity) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// OpenBot 开放平台三大能力之一：群通知机器人（Webhook 单向推消息，不是数字伙伴）。
//
// BotID = IM 用户；进群走普通群成员，本表不存群关系。
// Webhook 入站后推送到该 Bot 所在的全部群。
// 可见性与 Digital 一致：谁能搜到/拉进群由 Visibility（+ Admin / Viewer）控制。
// Token 仅存哈希（bearer 比对）；SignatureSecret 需明文用于 HMAC 重算验签，生产建议 KMS 加密存储；
// 明文只在创建或重置时返回一次。
type OpenBot struct {
	gorm.Model
	OwnerID    string `gorm:"column:owner_id;type:varchar(64);index;not null;default:'';comment:创建者用户ID"`
	BotID      string `gorm:"column:bot_id;type:varchar(64);uniqueIndex;not null;comment:Bot 的 IM 用户ID"`
	Token      string `gorm:"column:token_hash;type:varchar(128);uniqueIndex;not null;comment:Webhook Token 哈希"`
	Visibility string `gorm:"type:varchar(16);not null;default:private;index;comment:private|partial|public"`
	Status     int    `gorm:"column:status;type:tinyint;not null;default:1;index;comment:1启用 0禁用"`

	SignatureEnabled int    `gorm:"column:signature_enabled;type:tinyint;not null;default:0;comment:是否启用签名校验 1是 0否"`
	SignatureSecret  string `gorm:"column:signature_secret;type:varchar(256);not null;default:'';comment:签名密钥（HMAC 验签需明文，生产建议 KMS 加密存储）"`

	Security OpenBotSecurity `gorm:"type:json;comment:入站安全配置（关键词/IP白名单）"`
}

func (OpenBot) TableName() string { return "open_bot" }
