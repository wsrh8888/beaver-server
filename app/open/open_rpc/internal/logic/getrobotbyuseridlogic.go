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

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type GetRobotByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetRobotByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotByUserIDLogic {
	return &GetRobotByUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_robot_by_user_id", ctx),
	}
}

func (l *GetRobotByUserIDLogic) GetRobotByUserID(in *open_rpc.GetRobotByUserIDReq) (*open_rpc.GetRobotByUserIDRes, error) {
	if in.RobotUserId == "" {
		return &open_rpc.GetRobotByUserIDRes{Found: false}, nil
	}

	var robot open_models.OpenRobot
	if err := l.svcCtx.DB.Where("robot_user_id = ? AND status = 1", in.RobotUserId).First(&robot).Error; err != nil {
		return &open_rpc.GetRobotByUserIDRes{Found: false}, nil
	}

	// 单聊/群聊/@ 响应行为不再由平台存储，由机器人自己决定，这里不再返回开关。
	return &open_rpc.GetRobotByUserIDRes{
		Found:       true,
		AppId:       robot.AppID,
		RobotUserId: robot.RobotID,
	}, nil
}
