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
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	// 2. 根据能力类型更新对应的开关
	var enabled bool
	switch req.Capability {
	case "robot":
		if req.Enable {
			app.EnableRobot = 1
			app.EnableWebhook = 1
			enabled = true
		} else {
			app.EnableRobot = 0
			enabled = false
		}
	case "oauth":
		if req.Enable {
			app.EnableOAuth = 1
			enabled = true
		} else {
			app.EnableOAuth = 0
			enabled = false
		}
	case "webhook":
		if req.Enable {
			app.EnableWebhook = 1
			enabled = true
		} else {
			app.EnableWebhook = 0
			enabled = false
		}
	default:
		return nil, errors.New("不支持的能力类型")
	}

	// 3. 保存更新
	if err := l.svcCtx.DB.Save(&app).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "更新应用能力失败", Data: map[string]interface{}{"err": err}})
		return nil, errors.New("更新应用能力失败")
	}

	if req.Capability == "robot" && req.Enable {
		if err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app); err != nil {
			l.logger.Error(model.LogMsg{Text: "创建 Robot 用户失败", Data: map[string]interface{}{"app_id": req.AppID, "err": err.Error()}})
			return nil, errors.New("启用 Robot 成功，但创建 IM 用户失败，请稍后重试")
		}
	}

	l.logger.Info(model.LogMsg{Text: "应用能力开关已更新", Data: map[string]interface{}{"app_id": req.AppID, "capability": req.Capability, "enabled": req.Enable}})

	return &types.ToggleAppCapabilityRes{
		Enabled: enabled,
	}, nil
}

func ensurePortalAppRobot(ctx context.Context, db *gorm.DB, userRpc user.User, app *open_models.OpenApp) error {
	var robot open_models.OpenAppRobot
	err := db.Where("app_id = ?", app.AppID).First(&robot).Error
	if err == nil && robot.RobotID != "" {
		return nil
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

	robot = open_models.OpenAppRobot{
		AppID:            app.AppID,
		RobotID:          createRes.UserID,
		RobotName:        nickName,
		Avatar:           app.Icon,
		Status:           1,
		EnableSingleChat: 1,
		EnableGroupChat:  1,
		EnableAtMention:  1,
	}
	return db.Save(&robot).Error
}
