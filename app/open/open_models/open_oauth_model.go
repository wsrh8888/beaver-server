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

// OpenOAuth 开放平台三大能力之一：OAuth 快捷登录（凭据 + OAuth 配置同表，与 Bot/Digital 无关）。
// OAuthSecret 仅存单向哈希（或 KMS 密文），明文只在签发/重置时返回一次。
// 对外标识统一用 OAuthID。
type OpenOAuth struct {
	gorm.Model
	OAuthID     string `gorm:"column:oauth_id;type:varchar(64);uniqueIndex;not null;comment:OAuth 应用标识 oauth_xxx"`
	OAuthSecret string `gorm:"column:oauth_secret_hash;type:varchar(128);not null;comment:OAuth 密钥哈希"`
	Name        string `gorm:"type:varchar(100);not null;comment:应用名称"`
	Description string `gorm:"type:text;comment:应用描述"`
	Icon        string `gorm:"type:varchar(500);comment:应用图标URL"`
	OwnerID     string `gorm:"column:owner_id;type:varchar(64);index;comment:所属用户ID"`
	Status      int    `gorm:"type:tinyint;not null;default:1;index;comment:1启用 0禁用"`

	RequirePKCE     int    `gorm:"type:tinyint;not null;default:0;comment:是否强制 PKCE 1是 0否"`
	SupportedScopes string `gorm:"type:json;comment:支持的权限范围 JSON 数组"`
	AccessTokenTTL  int    `gorm:"type:int;not null;default:7200;comment:Access Token 有效期秒"`
	RefreshTokenTTL int    `gorm:"type:int;not null;default:2592000;comment:Refresh Token 有效期秒"`

	H5      *H5OAuth      `gorm:"type:json;comment:H5 OAuth 配置"`
	Desktop *DesktopOAuth `gorm:"type:json;comment:桌面端 OAuth 配置"`
	Mobile  *MobileOAuth  `gorm:"type:json;comment:移动端 OAuth 配置"`
}

func (OpenOAuth) TableName() string { return "open_oauth" }

type H5OAuth struct {
	Enabled      bool     `json:"enabled"`
	RedirectURIs []string `json:"redirectUris"`
}

func (c *H5OAuth) Value() (driver.Value, error) { return json.Marshal(c) }
func (c *H5OAuth) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

type DesktopOAuth struct {
	Enabled      bool   `json:"enabled"`
	CustomScheme string `json:"customScheme"`
}

func (c *DesktopOAuth) Value() (driver.Value, error) { return json.Marshal(c) }
func (c *DesktopOAuth) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}

type MobileOAuth struct {
	Enabled            bool   `json:"enabled"`
	IOSBundleID        string `json:"iosBundleId"`
	AndroidPackageName string `json:"androidPackageName"`
	UniversalLink      string `json:"universalLink"`
	CustomScheme       string `json:"customScheme"`
}

func (c *MobileOAuth) Value() (driver.Value, error) { return json.Marshal(c) }
func (c *MobileOAuth) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}
