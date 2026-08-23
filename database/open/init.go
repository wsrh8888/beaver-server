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

package openseed

import (
	"encoding/json"
	"fmt"
	"log"

	"beaver/app/open/constants"
	"beaver/app/open/open_models"

	"gorm.io/gorm"
)

const (
	// 与 beaver-desktop / beaver-open 的 openAppId 保持一致
	defaultAppID     = "app_2db39e38"
	defaultAppSecret = "beaver_im_oauth_secret_dev"
	defaultOwnerUserID = "100000"
	defaultAppName     = "海狸IM"
	defaultAppDesc     = "系统初始化默认应用（OAuth 登录）"
	defaultDesktopScheme = "beaver://"
)

// InitQuickLoginApp 初始化开放平台默认「海狸IM」应用（启用 Desktop + H5 OAuth，不含机器人等能力）
func InitQuickLoginApp(db *gorm.DB) error {
	var app open_models.OpenApp
	err := db.Where("app_id = ?", defaultAppID).First(&app).Error
	if err == gorm.ErrRecordNotFound {
		app = open_models.OpenApp{
			AppID:       defaultAppID,
			AppSecret:   defaultAppSecret,
			Name:        defaultAppName,
			Description: defaultAppDesc,
			OwnerUserID: defaultOwnerUserID,
			AppType:     1,
			Category:    "tool",
			Status:      1, // 已发布
			AuditStatus: 1, // 已通过
			EnableOAuth: 1,
			Scheme:      defaultDesktopScheme,
			Version:     1,
		}
		if err := db.Create(&app).Error; err != nil {
			return fmt.Errorf("创建默认应用失败: %w", err)
		}
		log.Printf("创建默认应用成功: appId=%s name=%s", defaultAppID, defaultAppName)
	} else if err != nil {
		return fmt.Errorf("查询默认应用失败: %w", err)
	} else {
		updates := map[string]interface{}{}
		if app.Name != defaultAppName {
			updates["name"] = defaultAppName
		}
		if app.Description != defaultAppDesc {
			updates["description"] = defaultAppDesc
		}
		if app.Status != 1 {
			updates["status"] = 1
		}
		if app.AuditStatus != 1 {
			updates["audit_status"] = 1
		}
		if app.EnableOAuth != 1 {
			updates["enable_oauth"] = 1
		}
		if app.Scheme == "" {
			updates["scheme"] = defaultDesktopScheme
		}
		if len(updates) > 0 {
			if err := db.Model(&app).Updates(updates).Error; err != nil {
				return fmt.Errorf("更新默认应用失败: %w", err)
			}
			log.Printf("默认应用已校正: appId=%s", defaultAppID)
		} else {
			log.Printf("默认应用已存在: appId=%s", defaultAppID)
		}
	}

	scopes, _ := json.Marshal([]string{
		string(constants.ScopeUserProfileRead),
		string(constants.ScopeUserAvatarRead),
	})

	desktopCfg := &open_models.DesktopOAuth{
		Enabled:      true,
		CustomScheme: defaultDesktopScheme,
	}
	h5Cfg := &open_models.H5OAuth{
		Enabled: true,
		// beaver-open 本地开发默认回调（hash 模式下 code 落在 origin/）
		RedirectURIs: []string{
			"http://localhost:4012/",
			"http://127.0.0.1:4012/",
			"https://fe.wsrh8888.com/open/",
		},
		JsSdkDomains: []string{
			"http://localhost:4012",
			"http://127.0.0.1:4012",
			"https://fe.wsrh8888.com",
		},
	}

	var oauth open_models.OpenAppOAuth
	err = db.Where("app_id = ?", defaultAppID).First(&oauth).Error
	if err == gorm.ErrRecordNotFound {
		oauth = open_models.OpenAppOAuth{
			AppID:           defaultAppID,
			SupportedScopes: string(scopes),
			AccessTokenTTL:  7200,
			RefreshTokenTTL: 2592000,
			Desktop:         desktopCfg,
			H5:              h5Cfg,
		}
		if err := db.Create(&oauth).Error; err != nil {
			return fmt.Errorf("创建默认 OAuth 配置失败: %w", err)
		}
		log.Printf("创建默认 OAuth 配置成功: appId=%s desktop+h5", defaultAppID)
		return nil
	}
	if err != nil {
		return fmt.Errorf("查询默认 OAuth 配置失败: %w", err)
	}

	needSave := false
	if oauth.Desktop == nil || !oauth.Desktop.Enabled {
		oauth.Desktop = desktopCfg
		needSave = true
	} else if oauth.Desktop.CustomScheme == "" {
		oauth.Desktop.CustomScheme = defaultDesktopScheme
		needSave = true
	}
	if oauth.H5 == nil || !oauth.H5.Enabled {
		oauth.H5 = h5Cfg
		needSave = true
	} else {
		if len(oauth.H5.RedirectURIs) == 0 {
			oauth.H5.RedirectURIs = h5Cfg.RedirectURIs
			needSave = true
		}
		if len(oauth.H5.JsSdkDomains) == 0 {
			oauth.H5.JsSdkDomains = h5Cfg.JsSdkDomains
			needSave = true
		}
	}
	if oauth.SupportedScopes == "" {
		oauth.SupportedScopes = string(scopes)
		needSave = true
	}
	if needSave {
		if err := db.Save(&oauth).Error; err != nil {
			return fmt.Errorf("更新默认 OAuth 配置失败: %w", err)
		}
		log.Printf("默认 OAuth 配置已校正: appId=%s desktop+h5", defaultAppID)
	} else {
		log.Printf("默认 OAuth 配置已存在: appId=%s", defaultAppID)
	}

	return nil
}
