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

package robot

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
	"gorm.io/gorm"
)

type GetRobotConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetRobotConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotConfigLogic {
	return &GetRobotConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_robot_config", ctx),
	}
}

func (l *GetRobotConfigLogic) GetRobotConfig(req *types.GetRobotConfigReq) (resp *types.GetRobotConfigRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	robot, err := ensurePortalAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, &app)
	if err != nil {
		return nil, errors.New("获取 Robot 配置失败")
	}

	config := types.RobotConfigInfo{
		AppID:   req.AppID,
		RobotID: robot.RobotID,
		Status:  robot.Status,
	}
	// 昵称/头像取 IM 用户资料（RobotID 即 IM 用户ID）。
	// 欢迎语/命令前缀/单聊群聊开关不再由平台存储，由机器人自己维护。
	if userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoReq{UserID: robot.RobotID}); err == nil && userRes.UserInfo != nil {
		config.RobotName = userRes.UserInfo.NickName
		config.Avatar = userRes.UserInfo.Avatar
	}

	return &types.GetRobotConfigRes{
		Config: config,
	}, nil
}

func ensurePortalAppRobot(ctx context.Context, db *gorm.DB, userRpc user.User, app *open_models.OpenApp) (*open_models.OpenRobot, error) {
	var robot open_models.OpenRobot
	err := db.Where("app_id = ?", app.AppID).First(&robot).Error
	if err == nil && robot.RobotID != "" {
		return &robot, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
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
		return nil, fmt.Errorf("user create: %w", err)
	}

	robot = open_models.OpenRobot{
		AppID:   app.AppID,
		RobotID: createRes.UserID,
		Status:  1,
	}
	if err := db.Save(&robot).Error; err != nil {
		return nil, err
	}
	return &robot, nil
}