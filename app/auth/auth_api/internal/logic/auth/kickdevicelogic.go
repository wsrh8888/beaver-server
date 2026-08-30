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
	"errors"
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/app/auth/auth_models"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type KickDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewKickDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickDeviceLogic {
	return &KickDeviceLogic{
		ctx:    ctx,
		logger: logger.New("kick_device"),
		svcCtx: svcCtx,
	}
}

func (l *KickDeviceLogic) KickDevice(req *types.KickDeviceReq) (*types.KickDeviceRes, error) {
	var device auth_models.AuthDeviceModel
	if err := l.svcCtx.DB.Where("user_id = ? AND device_id = ?", req.UserID, req.DeviceID).
		First(&device).Error; err != nil {
		return nil, errors.New("设备不存在")
	}

	redisKey := fmt.Sprintf("user_authentication_session:%s:%s", req.UserID, device.DeviceType)
	l.svcCtx.Redis.Del(redisKey)

	l.svcCtx.DB.Model(&device).Update("is_active", false)

	if _, err := l.svcCtx.NotificationRpc.DisablePushToken(l.ctx, &notification_rpc.DisablePushTokenReq{
		UserId:   req.UserID,
		DeviceId: req.DeviceID,
	}); err != nil {
		logx.WithContext(l.ctx).Errorf("禁用 Push Token 失败: userId=%s, deviceId=%s, err=%v", req.UserID, req.DeviceID, err)
	}

	go func() {
		payload := map[string]interface{}{
			"command":  wsCommandConst.USER_PROFILE,
			"type":     wsTypeConst.UserKickReceive,
			"senderId": req.UserID,
			"targetId": req.UserID,
			"body": map[string]interface{}{
				"deviceId": req.DeviceID,
				"reason":   "kicked_by_user",
			},
			"conversationId": "",
		}
		l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)
	}()

	l.logger.Info(model.LogMsg{
		Text: "用户踢出设备成功",
		Data: map[string]interface{}{
			"userId":   req.UserID,
			"deviceId": req.DeviceID,
		},
	})
	return &types.KickDeviceRes{}, nil
}
