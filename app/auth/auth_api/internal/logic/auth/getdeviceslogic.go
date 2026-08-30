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
	"context"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/core/coreonline"
	"beaver/utils/device"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDevicesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDevicesLogic {
	return &GetDevicesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDevicesLogic) GetDevices(req *types.GetDevicesReq) (*types.GetDevicesRes, error) {
	var devices []auth_models.AuthDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ?", req.UserID).Order("last_login_time DESC").Find(&devices).Error; err != nil {
		l.Errorf("查询设备失败: %v", err)
		return nil, err
	}

	onlineDeviceIDs := make(map[string]bool)
	for _, slot := range []string{"desktop", "mobile"} {
		if !coreonline.IsSlotOnline(l.svcCtx.Redis, req.UserID, slot) {
			continue
		}
		deviceID, err := device.SessionDeviceID(l.svcCtx.Redis, req.UserID, slot)
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		onlineDeviceIDs[deviceID] = true
	}

	list := make([]types.DeviceInfo, 0, len(devices))
	for _, d := range devices {
		if !d.IsActive {
			continue
		}
		list = append(list, types.DeviceInfo{
			DeviceID:        d.DeviceID,
			DeviceType:      d.DeviceType,
			DeviceOS:        d.DeviceOS,
			DeviceModel:     d.DeviceModel,
			DeviceOsVersion: d.DeviceOsVersion,
			DeviceName:      d.DeviceName,
			LastLoginTime:   d.LastLoginTime.Format("2006-01-02 15:04:05"),
			IsOnline:        onlineDeviceIDs[d.DeviceID],
			LastLoginIP:     d.LastLoginIP,
		})
	}
	return &types.GetDevicesRes{Devices: list}, nil
}
