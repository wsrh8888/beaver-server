/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package utils

import (
	"errors"
	"strings"
	"time"

	"beaver/app/open/open_models"

	"gorm.io/gorm"
)

func parseBearerToken(authorization string) string {
	if authorization == "" {
		return ""
	}
	parts := strings.SplitN(authorization, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return authorization
}

func ValidateAppAccessToken(db *gorm.DB, authorization string) (*open_models.OpenOAuthToken, error) {
	token := parseBearerToken(authorization)
	if token == "" {
		return nil, errors.New("缺少访问令牌")
	}

	var record open_models.OpenOAuthToken
	if err := db.Where("access_token_hash = ? AND revoked_at IS NULL", token).First(&record).Error; err != nil {
		return nil, errors.New("访问令牌无效")
	}
	if time.Now().After(record.ExpiresAt) {
		return nil, errors.New("访问令牌已过期")
	}
	return &record, nil
}

// ValidateDigitalAccessToken 校验外接数字伙伴短时 access_token。
func ValidateDigitalAccessToken(db *gorm.DB, authorization string) (*open_models.OpenDigitalToken, error) {
	token := parseBearerToken(authorization)
	if token == "" {
		return nil, errors.New("缺少访问令牌")
	}

	var record open_models.OpenDigitalToken
	if err := db.Where("access_token_hash = ? AND revoked_at IS NULL", token).First(&record).Error; err != nil {
		return nil, errors.New("访问令牌无效")
	}
	if time.Now().After(record.ExpiresAt) {
		return nil, errors.New("访问令牌已过期")
	}
	return &record, nil
}

// RequireOAuthEnabled 校验 OAuth 应用是否处于可用状态。
func RequireOAuthEnabled(app *open_models.OpenOAuth) error {
	if app == nil {
		return errors.New("应用不存在")
	}
	if app.Status != 1 {
		return errors.New("应用已禁用")
	}
	return nil
}

func LoadOAuthByID(db *gorm.DB, oauthID string) (*open_models.OpenOAuth, error) {
	var app open_models.OpenOAuth
	if err := db.Where("oauth_id = ?", oauthID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在")
	}
	return &app, nil
}
