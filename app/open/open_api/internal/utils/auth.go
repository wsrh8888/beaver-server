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
	if err := db.Where("token = ?", token).First(&record).Error; err != nil {
		return nil, errors.New("访问令牌无效")
	}
	if time.Now().Unix() > record.ExpiresAt {
		return nil, errors.New("访问令牌已过期")
	}
	return &record, nil
}

func RequireAppCapability(app *open_models.OpenApp, needRobot, needWebhook bool) error {
	if app.Status != 1 {
		return errors.New("应用未发布或已禁用")
	}
	if needRobot && app.EnableRobot != 1 {
		return errors.New("应用未启用智能机器人能力")
	}
	if needWebhook && app.EnableWebhook != 1 {
		return errors.New("应用未启用 Webhook 能力")
	}
	return nil
}

func LoadAppByID(db *gorm.DB, appID string) (*open_models.OpenApp, error) {
	var app open_models.OpenApp
	if err := db.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在")
	}
	return &app, nil
}
