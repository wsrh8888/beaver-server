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

	"beaver/app/call/call_models"
	"beaver/app/call/call_rpc/internal/svc"
	"beaver/app/call/call_rpc/types/call_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSessionLogic {
	return &GetSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话基础信息
func (l *GetSessionLogic) GetSession(in *call_rpc.GetSessionReq) (*call_rpc.GetSessionRes, error) {
	var session call_models.CallSession
	if err := l.svcCtx.DB.Where("room_id = ?", in.RoomId).First(&session).Error; err != nil {
		return nil, err
	}

	var participantIDs []string
	l.svcCtx.DB.Model(&call_models.CallParticipant{}).
		Where("room_id = ?", in.RoomId).
		Pluck("user_id", &participantIDs)

	return &call_rpc.GetSessionRes{
		RoomId:         session.RoomID,
		CallerId:       session.CallerID,
		CallType:       int32(session.CallType),
		Status:         int32(session.Status),
		ParticipantIds: participantIDs,
		ConversationId: session.ConversationID,
	}, nil
}
