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

package app

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/user"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"gorm.io/gorm"
)

type ToggleAppCapabilityLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 启用/禁用应用能力（对标飞书）
func NewToggleAppCapabilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ToggleAppCapabilityLogic {
	return &ToggleAppCapabilityLogic{
		logger: beaverlog.New("toggle_app_capability", ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToggleAppCapabilityLogic) ToggleAppCapability(req *types.ToggleAppCapabilityReq) (resp *types.ToggleAppCapabilityRes, err error) {

	// 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	// 2. 应用表不再持有能力开关：「具备某能力」由对应能力表是否存在有效记录表达。
	var enabled bool
	switch req.Capability {
	case "digital":
		if req.Enable {
			if _, err := ensurePortalAppExternalPartner(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app); err != nil {
				l.logger.Error(model.LogMsg{Text: "创建外接数字伙伴失败", Data: map[string]interface{}{"app_id": req.AppID, "err": err.Error()}})
				return nil, errors.New("启用失败：创建数字伙伴 IM 用户失败，请稍后重试")
			}
			enabled = true
		} else {
			if err := l.svcCtx.DB.Model(&open_models.OpenDigitalPartner{}).
				Where("app_id = ? AND runtime = ?", app.AppID, open_models.PartnerRuntimeExternal).
				Update("status", 2).Error; err != nil {
				l.logger.Error(model.LogMsg{Text: "停用外接数字伙伴失败", Data: map[string]interface{}{"app_id": req.AppID, "err": err.Error()}})
				return nil, errors.New("停用失败")
			}
			enabled = false
		}
	case "oauth":
		if req.Enable {
			var oauthConfig open_models.OpenOAuthConfig
			err := l.svcCtx.DB.Where("app_id = ?", app.AppID).First(&oauthConfig).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := l.svcCtx.DB.Create(&open_models.OpenOAuthConfig{AppID: app.AppID}).Error; err != nil {
					l.logger.Error(model.LogMsg{Text: "启用 OAuth 失败", Data: map[string]interface{}{"app_id": req.AppID, "err": err.Error()}})
					return nil, errors.New("启用 OAuth 失败")
				}
			} else if err != nil {
				return nil, errors.New("启用 OAuth 失败")
			}
			enabled = true
		} else {
			if err := l.svcCtx.DB.Unscoped().Where("app_id = ?", app.AppID).
				Delete(&open_models.OpenOAuthConfig{}).Error; err != nil {
				l.logger.Error(model.LogMsg{Text: "停用 OAuth 失败", Data: map[string]interface{}{"app_id": req.AppID, "err": err.Error()}})
				return nil, errors.New("停用 OAuth 失败")
			}
			enabled = false
		}
	case "webhook":
		return nil, errors.New("Webhook 能力已下线，平台事件改由长连接下发")
	default:
		return nil, errors.New("不支持的能力类型")
	}

	l.logger.Info(model.LogMsg{Text: "应用能力已更新", Data: map[string]interface{}{"app_id": req.AppID, "capability": req.Capability, "enabled": req.Enable}})

	return &types.ToggleAppCapabilityRes{
		Enabled: enabled,
	}, nil
}

// ensurePortalAppExternalPartner 确保应用挂有外接数字伙伴（PartnerID = IM 用户 ID）。
func ensurePortalAppExternalPartner(ctx context.Context, db *gorm.DB, userRpc user.User, app *open_models.OpenApp) (*open_models.OpenDigitalPartner, error) {
	var partner open_models.OpenDigitalPartner
	err := db.Where(
		"app_id = ? AND runtime = ?",
		app.AppID, open_models.PartnerRuntimeExternal,
	).First(&partner).Error
	if err == nil && partner.PartnerID != "" {
		if partner.Status != 1 {
			if err := db.Model(&partner).Update("status", 1).Error; err != nil {
				return nil, err
			}
			partner.Status = 1
		}
		return &partner, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	nickName := app.Name
	if nickName == "" {
		nickName = "数字伙伴"
	}

	createRes, err := userRpc.UserCreate(ctx, &user.UserCreateReq{
		NickName: nickName,
		UserType: int32(user_models.UserTypeDigital),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return nil, fmt.Errorf("user create: %w", err)
	}

	partner = open_models.OpenDigitalPartner{
		PartnerID:  createRes.UserID,
		OwnerID:    app.OwnerID,
		Runtime:    open_models.PartnerRuntimeExternal,
		AppID:      app.AppID,
		Visibility: open_models.PartnerVisibilityPrivate,
		Status:     1,
	}
	if err := db.Create(&partner).Error; err != nil {
		return nil, err
	}
	return &partner, nil
}
