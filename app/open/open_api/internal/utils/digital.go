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

	"beaver/app/open/open_models"

	"gorm.io/gorm"
)

func LoadDigitalByAppKey(db *gorm.DB, appKey string) (*open_models.OpenDigital, error) {
	if appKey == "" {
		return nil, errors.New("缺少 appKey")
	}
	var partner open_models.OpenDigital
	if err := db.Where(
		"app_key = ? AND runtime = ?",
		appKey, open_models.DigitalRuntimeExternal,
	).First(&partner).Error; err != nil {
		return nil, errors.New("数字伙伴不存在")
	}
	if partner.DigitalID == "" {
		return nil, errors.New("数字伙伴身份未就绪")
	}
	return &partner, nil
}

func LoadDigitalByDigitalID(db *gorm.DB, partnerID string) (*open_models.OpenDigital, error) {
	if partnerID == "" {
		return nil, errors.New("缺少 digitalId")
	}
	var partner open_models.OpenDigital
	if err := db.Where("digital_id = ?", partnerID).First(&partner).Error; err != nil {
		return nil, errors.New("数字伙伴不存在")
	}
	return &partner, nil
}

func RequireDigitalEnabled(partner *open_models.OpenDigital) error {
	if partner == nil {
		return errors.New("数字伙伴不存在")
	}
	if partner.Status != open_models.DigitalStatusEnabled {
		return errors.New("数字伙伴已禁用")
	}
	return nil
}

// ResolveDigitalFromToken 从 access token 的 DigitalID 解析外接数字伙伴。
func ResolveDigitalFromToken(db *gorm.DB, digitalID string) (*open_models.OpenDigital, error) {
	partner, err := LoadDigitalByDigitalID(db, digitalID)
	if err != nil {
		return nil, err
	}
	if err := RequireDigitalEnabled(partner); err != nil {
		return nil, err
	}
	return partner, nil
}
