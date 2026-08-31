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

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/core/coreonline"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"beaver/utils/device"

	"github.com/go-redis/redis"
)

type ListUserDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListUserDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserDevicesLogic {
	return &ListUserDevicesLogic{
		ctx:    ctx,
		logger: beaverlog.New("list_user_devices", ctx),
		svcCtx: svcCtx,
	}
}

func (l *ListUserDevicesLogic) ListUserDevices(in *auth_rpc.ListUserDevicesReq) (*auth_rpc.ListUserDevicesRes, error) {
	var devices []auth_models.AuthDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ? AND is_active = ?", in.UserId, true).
		Order("last_login_time DESC").
		Find(&devices).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询用户设备失败",
			Data: map[string]any{"userId": in.UserId, "err": err.Error()},
		})
		return nil, err
	}

	onlineDeviceIDs := make(map[string]bool)
	for _, slot := range []string{"desktop", "mobile"} {
		if !coreonline.IsSlotOnline(l.svcCtx.Redis, in.UserId, slot) {
			continue
		}
		deviceID, err := device.SessionDeviceID(l.svcCtx.Redis, in.UserId, slot)
		if err == redis.Nil {
			continue
		}
		if err != nil {
			l.logger.Error(model.LogMsg{
				Text: "读取会话设备失败",
				Data: map[string]any{
					"userId": in.UserId,
					"slot":   slot,
					"err":    err.Error(),
				},
			})
			return nil, err
		}
		onlineDeviceIDs[deviceID] = true
	}

	list := make([]*auth_rpc.UserDeviceItem, 0, len(devices))
	for _, d := range devices {
		list = append(list, &auth_rpc.UserDeviceItem{
			DeviceId:        d.DeviceID,
			DeviceType:      d.DeviceType,
			DeviceOs:        d.DeviceOS,
			DeviceModel:     d.DeviceModel,
			DeviceOsVersion: d.DeviceOsVersion,
			DeviceName:      d.DeviceName,
			LastLoginTime:   d.LastLoginTime.Format("2006-01-02 15:04:05"),
			LastLoginIp:     d.LastLoginIP,
			IsOnline:        onlineDeviceIDs[d.DeviceID],
		})
	}

	l.logger.Info(model.LogMsg{
		Text: "查询用户设备成功",
		Data: map[string]any{
			"userId":      in.UserId,
			"count":       len(list),
			"onlineCount": len(onlineDeviceIDs),
		},
	})
	return &auth_rpc.ListUserDevicesRes{List: list}, nil
}
