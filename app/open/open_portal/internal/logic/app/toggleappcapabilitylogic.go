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
	case "robot":
		if req.Enable {
			if err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app); err != nil {
				l.logger.Error(model.LogMsg{Text: "创建 Robot 用户失败", Data: map[string]interface{}{"app_id": req.AppID, "err": err.Error()}})
				return nil, errors.New("启用 Robot 失败：创建 IM 用户失败，请稍后重试")
			}
			enabled = true
		} else {
			if err := l.svcCtx.DB.Model(&open_models.OpenRobot{}).
				Where("app_id = ?", app.AppID).Update("status", 0).Error; err != nil {
				l.logger.Error(model.LogMsg{Text: "停用 Robot 失败", Data: map[string]interface{}{"app_id": req.AppID, "err": err.Error()}})
				return nil, errors.New("停用 Robot 失败")
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

// ensurePortalAppRobot 确保应用已有可用的 Robot 记录（RobotID 即 IM 用户ID）。
// 昵称/头像不再由平台冗余存储，行为配置由机器人自己维护。
func ensurePortalAppRobot(ctx context.Context, db *gorm.DB, userRpc user.User, app *open_models.OpenApp) error {
	var robot open_models.OpenRobot
	err := db.Where("app_id = ?", app.AppID).First(&robot).Error
	if err == nil && robot.RobotID != "" {
		if robot.Status == 1 {
			return nil
		}
		return db.Model(&robot).Update("status", 1).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	nickName := app.Name
	if nickName == "" {
		nickName = "Robot"
	}

	createRes, err := userRpc.UserCreate(ctx, &user.UserCreateReq{
		NickName: nickName,
		UserType: int32(user_models.UserTypeRobot),
		Source:   int32(user_models.SourceGroup),
	})
	if err != nil {
		return fmt.Errorf("user create: %w", err)
	}

	robot = open_models.OpenRobot{
		AppID:   app.AppID,
		RobotID: createRes.UserID,
		Status:  1,
	}
	return db.Save(&robot).Error
}
