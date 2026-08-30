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
	"gorm.io/gorm"
)

type CreateSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSessionLogic {
	return &CreateSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 核心：创建通话会话并初始化参与者名单
func (l *CreateSessionLogic) CreateSession(in *call_rpc.CreateSessionReq) (*call_rpc.CreateSessionRes, error) {
	var participants []*call_rpc.Participant

	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建会话
		session := &call_models.CallSession{
			RoomID:         in.RoomId,
			CallerID:       in.CallerId,
			CallType:       int8(in.CallType),
			ConversationID: in.ConversationId,                // 存入会话ID
			Status:         call_models.SessionStatusCalling, // 1-进行中
		}
		if err := tx.Create(session).Error; err != nil {
			return err
		}

		// 2. 创建参与者 (发起者)
		caller := &call_models.CallParticipant{
			RoomID: in.RoomId,
			UserID: in.CallerId,
			Status: call_models.ParticipantStatusJoined, // 2-已接听 (发起者默认已进入)
			Role:   1,                                   // 1-发起者
		}
		if err := tx.Create(caller).Error; err != nil {
			return err
		}
		participants = append(participants, &call_rpc.Participant{
			UserId: in.CallerId,
			Status: int32(call_models.ParticipantStatusJoined),
		})

		// 3. 创建参与者 (受邀者)
		if in.CallType == 1 {
			// 单聊：初始化对方
			target := &call_models.CallParticipant{
				RoomID: in.RoomId,
				UserID: in.TargetId,
				Status: call_models.ParticipantStatusCalling, // 1-待接听
				Role:   2,                                    // 2-受邀者
			}
			if err := tx.Create(target).Error; err != nil {
				return err
			}
			participants = append(participants, &call_rpc.Participant{
				UserId: in.TargetId,
				Status: int32(call_models.ParticipantStatusCalling),
			})
		} else {
			// 群聊：目前 TargetId 是群 ID，初始参与者仅发起人。
			// 后续参与者通过 Invite 会话接口加入，此时 GetParticipants 会返回最新名单。
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &call_rpc.CreateSessionRes{
		Success:      true,
		Participants: participants,
	}, nil
}
