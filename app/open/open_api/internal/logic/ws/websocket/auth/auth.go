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

package auth

import (
	"errors"

	"beaver/app/open/open_api/internal/utils"
	"beaver/app/open/open_models"

	"gorm.io/gorm"
)

// DigitalIdentity 外接数字伙伴建连鉴权结果。
type DigitalIdentity struct {
	AppID     string // 应用 ID
	DigitalID string // 数字伙伴 PartnerID = IM 用户 ID
}

// VerifyDigitalToken 校验建连凭据：ticket → app → 外接数字伙伴。
func VerifyDigitalToken(db *gorm.DB, appID, ticket string) (*DigitalIdentity, error) {
	if appID == "" {
		return nil, errors.New("缺少 appId")
	}
	if ticket == "" {
		return nil, errors.New("缺少 ticket")
	}

	token, err := utils.ValidateAppAccessToken(db, ticket)
	if err != nil {
		return nil, err
	}
	if token.AppID != appID {
		return nil, errors.New("ticket 与应用不匹配")
	}

	app, err := utils.LoadAppByID(db, appID)
	if err != nil {
		return nil, err
	}
	if err = utils.RequireAppEnabled(app); err != nil {
		return nil, err
	}

	var partner open_models.OpenDigitalPartner
	if err = db.Where(
		"app_id = ? AND runtime = ? AND status = 1",
		appID, open_models.PartnerRuntimeExternal,
	).First(&partner).Error; err != nil {
		return nil, errors.New("该应用未挂接外接数字伙伴")
	}
	if partner.PartnerID == "" {
		return nil, errors.New("数字伙伴身份未就绪")
	}

	return &DigitalIdentity{AppID: appID, DigitalID: partner.PartnerID}, nil
}
