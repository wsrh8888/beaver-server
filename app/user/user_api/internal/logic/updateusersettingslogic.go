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

package logic

import (
	"context"

	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"gorm.io/gorm"
)

type UpdateUserSettingsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdateUserSettingsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserSettingsLogic {
	return &UpdateUserSettingsLogic{
		ctx:    ctx,
		logger: logger.New("update_user_settings"),
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserSettingsLogic) UpdateUserSettings(req *types.UpdateUserSettingsReq) (*types.UpdateUserSettingsRes, error) {
	setting, err := l.getOrCreateUserSetting(req.UserID)
	if err != nil {
		return nil, err
	}

	defaults := user_models.DefaultUserSetting(req.UserID).SettingInfo
	if setting.SettingInfo == nil {
		setting.SettingInfo = defaults
	}
	if setting.SettingInfo.Privacy == nil {
		setting.SettingInfo.Privacy = defaults.Privacy
	}
	if setting.SettingInfo.Notification == nil {
		setting.SettingInfo.Notification = defaults.Notification
	}
	if setting.SettingInfo.Keyboard == nil {
		setting.SettingInfo.Keyboard = defaults.Keyboard
	}

	if req.Privacy != nil {
		if req.Privacy.AllowFriendRequest != nil {
			setting.SettingInfo.Privacy.AllowFriendRequest = *req.Privacy.AllowFriendRequest
		}
		if req.Privacy.ShowOnlineStatus != nil {
			setting.SettingInfo.Privacy.ShowOnlineStatus = *req.Privacy.ShowOnlineStatus
		}
		if req.Privacy.AllowSearchByPhone != nil {
			setting.SettingInfo.Privacy.AllowSearchByPhone = *req.Privacy.AllowSearchByPhone
		}
		if req.Privacy.AllowSearchByEmail != nil {
			setting.SettingInfo.Privacy.AllowSearchByEmail = *req.Privacy.AllowSearchByEmail
		}
	}

	if req.Notification != nil {
		if req.Notification.NotifyFriendRequest != nil {
			setting.SettingInfo.Notification.NotifyFriendRequest = *req.Notification.NotifyFriendRequest
		}
		if req.Notification.NotifyGroupMessage != nil {
			setting.SettingInfo.Notification.NotifyGroupMessage = *req.Notification.NotifyGroupMessage
		}
		if req.Notification.NotifyMoment != nil {
			setting.SettingInfo.Notification.NotifyMoment = *req.Notification.NotifyMoment
		}
	}

	if req.Keyboard != nil {
		if req.Keyboard.Screenshot != nil {
			setting.SettingInfo.Keyboard.Screenshot = *req.Keyboard.Screenshot
		}
		if req.Keyboard.ToggleWindow != nil {
			setting.SettingInfo.Keyboard.ToggleWindow = *req.Keyboard.ToggleWindow
		}
		if req.Keyboard.SendMessage != nil {
			setting.SettingInfo.Keyboard.SendMessage = *req.Keyboard.SendMessage
		}
	}

	if err := l.svcCtx.DB.Save(setting).Error; err != nil {
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "用户设置更新成功",
		Data: map[string]interface{}{"userId": req.UserID},
	})
	return &types.UpdateUserSettingsRes{}, nil
}

func (l *UpdateUserSettingsLogic) getOrCreateUserSetting(userID string) (*user_models.UserSettingModel, error) {
	var setting user_models.UserSettingModel
	err := l.svcCtx.DB.Where("user_id = ?", userID).First(&setting).Error
	if err == nil {
		return &setting, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	setting = user_models.DefaultUserSetting(userID)
	if err := l.svcCtx.DB.Create(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}
