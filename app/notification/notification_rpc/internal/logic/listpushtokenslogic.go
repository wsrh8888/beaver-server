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

type ListPushTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPushTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPushTokensLogic {
	return &ListPushTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPushTokensLogic) ListPushTokens(in *notification_rpc.ListPushTokensReq) (*notification_rpc.ListPushTokensRes, error) {
	if in.UserId == "" {
		return nil, errors.New("userId 不能为空")
	}

	var rows []notification_models.PushRegistrationModel
	if err := l.svcCtx.DB.Where("user_id = ? AND enabled = ?", in.UserId, true).Find(&rows).Error; err != nil {
		l.Errorf("查询 Push Token 失败: userId=%s, err=%v", in.UserId, err)
		return nil, errors.New("查询 Push Token 失败")
	}

	tokens := make([]*notification_rpc.PushTokenInfo, 0, len(rows))
	for _, row := range rows {
		tokens = append(tokens, &notification_rpc.PushTokenInfo{
			DeviceId:     row.DeviceID,
			PushToken:    row.PushToken,
			PushPlatform: row.PushPlatform,
		})
	}
	return &notification_rpc.ListPushTokensRes{Tokens: tokens}, nil
}
