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

package update_public

import (
	"context"
	"fmt"

	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/app/platform/platform_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type ReportVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewReportVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportVersionLogic {
	return &ReportVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("report_version", ctx),
	}
}

func (l *ReportVersionLogic) ReportVersion(req *types.ReportVersionReq) (*types.ReportVersionRes, error) {
	if req.ArchID == 0 {
		l.logger.Info(model.LogMsg{Text: "H5版本上报", Data: map[string]interface{}{"appId": req.AppID, "version": req.Version, "deviceId": req.DeviceID}})
		return &types.ReportVersionRes{}, nil
	}

	var app platform_models.UpdateApp
	if err := l.svcCtx.DB.Where("app_id = ? AND is_active = ?", req.AppID, true).First(&app).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("应用不存在或已停用")
		}
		l.logger.Error(model.LogMsg{Text: "查询应用失败", Data: map[string]any{"appId": req.AppID, "err": err.Error()}})
		return nil, fmt.Errorf("查询应用失败")
	}

	var architecture platform_models.UpdateArchitecture
	if err := l.svcCtx.DB.Where("app_id = ? AND platform_id = ? AND arch_id = ? AND is_active = ?",
		req.AppID, req.PlatformID, req.ArchID, true).First(&architecture).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("不支持的平台或架构")
		}
		l.logger.Error(model.LogMsg{Text: "查询架构失败", Data: map[string]any{"appId": req.AppID, "err": err.Error()}})
		return nil, fmt.Errorf("查询架构失败")
	}

	var existingReport platform_models.UpdateReport
	err := l.svcCtx.DB.Where("device_id = ? AND app_id = ? AND architecture_id = ?",
		req.DeviceID, req.AppID, architecture.Id).First(&existingReport).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		l.logger.Error(model.LogMsg{Text: "查询设备记录失败", Data: map[string]any{"deviceId": req.DeviceID, "appId": req.AppID, "err": err.Error()}})
		return nil, fmt.Errorf("查询设备记录失败")
	}

	if err == gorm.ErrRecordNotFound {
		report := platform_models.UpdateReport{
			UserID:         req.UserID,
			DeviceID:       req.DeviceID,
			AppID:          req.AppID,
			ArchitectureID: architecture.Id,
			Version:        req.Version,
		}
		if err := l.svcCtx.DB.Create(&report).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "创建版本上报记录失败", Data: map[string]any{"deviceId": req.DeviceID, "appId": req.AppID, "err": err.Error()}})
			return nil, fmt.Errorf("创建上报记录失败")
		}
	} else {
		updates := map[string]interface{}{
			"user_id": req.UserID,
			"version": req.Version,
		}
		if err := l.svcCtx.DB.Model(&existingReport).Updates(updates).Error; err != nil {
			l.logger.Error(model.LogMsg{Text: "更新版本上报记录失败", Data: map[string]any{"deviceId": req.DeviceID, "appId": req.AppID, "err": err.Error()}})
			return nil, fmt.Errorf("更新上报记录失败")
		}
	}

	l.logger.Info(model.LogMsg{
		Text: "版本上报成功",
		Data: map[string]interface{}{
			"appId":      req.AppID,
			"platformId": req.PlatformID,
			"archId":     req.ArchID,
			"version":    req.Version,
			"deviceId":   req.DeviceID,
		},
	})

	return &types.ReportVersionRes{}, nil
}
