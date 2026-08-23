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
	"fmt"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils/device"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		logger: logger.New("logout"),
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (*types.LogoutRes, error) {
	for _, group := range []string{"desktop", "mobile"} {
		key := fmt.Sprintf("user_authentication_session:%s:%s", req.UserID, group)
		l.svcCtx.Redis.Del(key)
	}

	if err := device.Deactivate(l.svcCtx.DB, req.UserID, req.DeviceID); err != nil {
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "用户登出成功",
		Data: map[string]interface{}{"userId": req.UserID, "deviceId": req.DeviceID},
	})
	return &types.LogoutRes{}, nil
}
