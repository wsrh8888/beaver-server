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

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CheckDeviceActiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckDeviceActiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckDeviceActiveLogic {
	return &CheckDeviceActiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CheckDeviceActiveLogic) CheckDeviceActive(in *auth_rpc.CheckDeviceActiveReq) (*auth_rpc.CheckDeviceActiveRes, error) {
	if in.UserId == "" || in.DeviceId == "" {
		return nil, errors.New("userId 和 deviceId 不能为空")
	}

	var device auth_models.AuthDeviceModel
	err := l.svcCtx.DB.Where("user_id = ? AND device_id = ? AND is_active = ?",
		in.UserId, in.DeviceId, true).First(&device).Error
	if err == gorm.ErrRecordNotFound {
		return &auth_rpc.CheckDeviceActiveRes{Active: false}, nil
	}
	if err != nil {
		l.Errorf("校验设备状态失败: userId=%s, deviceId=%s, err=%v", in.UserId, in.DeviceId, err)
		return nil, errors.New("校验设备状态失败")
	}

	return &auth_rpc.CheckDeviceActiveRes{Active: true}, nil
}
