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

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserStatusLogic {
	return &GetUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 供 Chat/User 等服务查询用户是否忙碌
func (l *GetUserStatusLogic) GetUserStatus(in *call_rpc.GetUserStatusReq) (*call_rpc.GetUserStatusRes, error) {
	var participant call_models.CallParticipant
	err := l.svcCtx.DB.Where("user_id = ? AND status IN ?", in.UserId, []int8{1, 2}).First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &call_rpc.GetUserStatusRes{
				IsBusy: false,
			}, nil
		}
		return nil, err
	}

	return &call_rpc.GetUserStatusRes{
		IsBusy: true,
		RoomId: participant.RoomID,
	}, nil
}
