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

// RobotIdentity 机器人建连鉴权结果
type RobotIdentity struct {
	AppID   string // 应用 ID（接入方身份）
	RobotID string // 机器人 IM 用户 ID（事件路由目标）
}

// VerifyRobotToken 校验机器人建连凭据，返回该连接代表的机器人身份。
//
// 凭据：appId + ticket。ticket 复用开放平台 OAPI 的 accessToken
// （POST /api/open/auth_public/v1/token 获取），因此服务端无需新增鉴权接口。
//
// 校验链：ticket 有效 → ticket 归属该 appId → 应用处于启用状态 → 应用已开启机器人能力。
//
// 与用户侧 ws_api 的 VerifyWsToken 完全分离：机器人没有用户登录态，
// 那条链路依赖 Redis 的 user_authentication_session:*，不可复用。
func VerifyRobotToken(db *gorm.DB, appID, ticket string) (*RobotIdentity, error) {
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

	var robot open_models.OpenRobot
	if err = db.Where("app_id = ? AND status = 1", appID).First(&robot).Error; err != nil {
		return nil, errors.New("该应用未开启机器人能力")
	}
	if robot.RobotID == "" {
		return nil, errors.New("机器人身份未就绪")
	}

	return &RobotIdentity{AppID: appID, RobotID: robot.RobotID}, nil
}
