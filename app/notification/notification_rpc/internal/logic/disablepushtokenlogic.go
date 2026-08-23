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
	"errors"

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DisablePushTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDisablePushTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisablePushTokenLogic {
	return &DisablePushTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DisablePushTokenLogic) DisablePushToken(in *notification_rpc.DisablePushTokenReq) (*notification_rpc.DisablePushTokenRes, error) {
	if in.UserId == "" || in.DeviceId == "" {
		return nil, errors.New("userId 和 deviceId 不能为空")
	}

	result := l.svcCtx.DB.Model(&notification_models.PushRegistrationModel{}).
		Where("user_id = ? AND device_id = ?", in.UserId, in.DeviceId).
		Update("enabled", false)
	if result.Error != nil {
		l.Errorf("禁用 Push Token 失败: userId=%s, deviceId=%s, err=%v", in.UserId, in.DeviceId, result.Error)
		return nil, errors.New("禁用 Push Token 失败")
	}

	return &notification_rpc.DisablePushTokenRes{Success: true}, nil
}
