/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package auth

import (
	"errors"

	"beaver/app/open/open_api/internal/utils"

	"gorm.io/gorm"
)

// DigitalIdentity 外接数字伙伴建连鉴权结果。
type DigitalIdentity struct {
	AppKey    string // 外接凭据 AppKey
	DigitalID string // DigitalID = IM 用户 ID
}

// VerifyDigitalToken 校验建连凭据：ticket → OpenDigitalToken.DigitalID → 数字伙伴，再校验 AppKey 匹配。
func VerifyDigitalToken(db *gorm.DB, appKey, ticket string) (*DigitalIdentity, error) {
	if appKey == "" {
		return nil, errors.New("缺少 appKey")
	}
	if ticket == "" {
		return nil, errors.New("缺少 ticket")
	}

	token, err := utils.ValidateDigitalAccessToken(db, ticket)
	if err != nil {
		return nil, err
	}

	// 令牌关联一级主键 DigitalID，回查伙伴后校验 AppKey 与建连凭据匹配。
	partner, err := utils.LoadDigitalByDigitalID(db, token.DigitalID)
	if err != nil {
		return nil, err
	}
	if partner.AppKey != appKey {
		return nil, errors.New("ticket 与 appKey 不匹配")
	}
	if err := utils.RequireDigitalEnabled(partner); err != nil {
		return nil, err
	}

	return &DigitalIdentity{AppKey: appKey, DigitalID: partner.DigitalID}, nil
}
