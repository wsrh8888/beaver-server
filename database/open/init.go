/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
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
	defaultOAuthID       = "oauth_2db39e38"
	defaultOAuthSecret   = "beaver_im_oauth_secret_dev"
	defaultOwnerUserID   = "100000"
	defaultAppName       = "海狸IM"
	defaultAppDesc       = "系统初始化默认应用（OAuth 登录）"
	defaultDesktopScheme = "beaver://"
)

// InitQuickLoginApp 初始化开放平台默认「海狸IM」OAuth 应用（凭据 + OAuth 配置同表）
func InitQuickLoginApp(db *gorm.DB) error {
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
		RedirectURIs: []string{
			"http://localhost:4012/",
			"http://127.0.0.1:4012/",
			"https://fe.wsrh8888.com/open/",
		},
	}

	var app open_models.OpenOAuth
	err := db.Where("oauth_id = ?", defaultOAuthID).First(&app).Error
	if err == gorm.ErrRecordNotFound {
		app = open_models.OpenOAuth{
			OAuthID:         defaultOAuthID,
			OAuthSecret:     defaultOAuthSecret,
			Name:            defaultAppName,
			Description:     defaultAppDesc,
			OwnerID:         defaultOwnerUserID,
			Status:          1,
			SupportedScopes: string(scopes),
			AccessTokenTTL:  7200,
			RefreshTokenTTL: 2592000,
			Desktop:         desktopCfg,
			H5:              h5Cfg,
		}
		if err := db.Create(&app).Error; err != nil {
			return fmt.Errorf("创建默认 OAuth 应用失败: %w", err)
		}
		log.Printf("创建默认 OAuth 应用成功: oauthId=%s desktop+h5", defaultOAuthID)
		return nil
	}
	if err != nil {
		return fmt.Errorf("查询默认 OAuth 应用失败: %w", err)
	}

	needSave := false
	if app.Name != defaultAppName {
		app.Name = defaultAppName
		needSave = true
	}
	if app.Description != defaultAppDesc {
		app.Description = defaultAppDesc
		needSave = true
	}
	if app.Status != 1 {
		app.Status = 1
		needSave = true
	}
	if app.SupportedScopes == "" {
		app.SupportedScopes = string(scopes)
		needSave = true
	}
	if app.AccessTokenTTL == 0 {
		app.AccessTokenTTL = 7200
		needSave = true
	}
	if app.RefreshTokenTTL == 0 {
		app.RefreshTokenTTL = 2592000
		needSave = true
	}
	if app.Desktop == nil || !app.Desktop.Enabled {
		app.Desktop = desktopCfg
		needSave = true
	} else if app.Desktop.CustomScheme == "" {
		app.Desktop.CustomScheme = defaultDesktopScheme
		needSave = true
	}
	if app.H5 == nil || !app.H5.Enabled {
		app.H5 = h5Cfg
		needSave = true
	} else if len(app.H5.RedirectURIs) == 0 {
		app.H5.RedirectURIs = h5Cfg.RedirectURIs
		needSave = true
	}
	if needSave {
		if err := db.Save(&app).Error; err != nil {
			return fmt.Errorf("更新默认 OAuth 应用失败: %w", err)
		}
		log.Printf("默认 OAuth 应用已校正: oauthId=%s desktop+h5", defaultOAuthID)
	} else {
		log.Printf("默认 OAuth 应用已存在: oauthId=%s", defaultOAuthID)
	}
	return nil
}
