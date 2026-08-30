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

package user_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PrivacySetting 隐私设置
type PrivacySetting struct {
	AllowFriendRequest bool `json:"allowFriendRequest"` // 是否允许他人发送好友申请
	ShowOnlineStatus   bool `json:"showOnlineStatus"`   // 是否向好友展示在线状态
	AllowSearchByPhone bool `json:"allowSearchByPhone"` // 是否允许通过手机号搜索到自己
	AllowSearchByEmail bool `json:"allowSearchByEmail"` // 是否允许通过邮箱搜索到自己
}

// NotificationSetting 通知策略
type NotificationSetting struct {
	NotifyFriendRequest bool `json:"notifyFriendRequest"` // 是否接收好友申请通知
	NotifyGroupMessage  bool `json:"notifyGroupMessage"`  // 是否接收群消息通知
	NotifyMoment        bool `json:"notifyMoment"`        // 是否接收朋友圈互动通知
}

// KeyboardSetting 快捷键（多端同步，键位字符串如 Ctrl+Alt+A）
type KeyboardSetting struct {
	Screenshot   string `json:"screenshot"`   // 区域截图
	ToggleWindow string `json:"toggleWindow"` // 显示/隐藏主窗口
	SendMessage  string `json:"sendMessage"`  // 发送消息（聊天输入框内生效）
}

// SettingInfo 用户设置详情（一人一条，JSON 内包含各模块）
type SettingInfo struct {
	Privacy      *PrivacySetting      `json:"privacy"`      // 隐私设置
	Notification *NotificationSetting `json:"notification"` // 通知策略
	Keyboard     *KeyboardSetting     `json:"keyboard"`     // 快捷键
}

func (s *SettingInfo) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SettingInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, s)
}

// UserSettingModel 用户设置，一个用户一行
type UserSettingModel struct {
	models.Model
	UserID      string       `gorm:"size:64;uniqueIndex" json:"userId"` // 用户 ID
	SettingInfo *SettingInfo `gorm:"type:longtext" json:"settingInfo"`  // 设置详情（JSON）
}

func (UserSettingModel) TableName() string {
	return "user_settings"
}

func DefaultUserSetting(userID string) UserSettingModel {
	return UserSettingModel{
		UserID: userID,
		SettingInfo: &SettingInfo{
			Privacy: &PrivacySetting{
				AllowFriendRequest: true,
				ShowOnlineStatus:   true,
				AllowSearchByPhone: true,
				AllowSearchByEmail: true,
			},
			Notification: &NotificationSetting{
				NotifyFriendRequest: true,
				NotifyGroupMessage:  true,
				NotifyMoment:        true,
			},
			Keyboard: &KeyboardSetting{
				Screenshot:   "Ctrl+Alt+A",
				ToggleWindow: "Ctrl+Shift+H",
				SendMessage:  "Enter",
			},
		},
	}
}
