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

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationLogic {
	return &UpdateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateConversationLogic) UpdateConversation(in *chat_rpc.UpdateConversationReq) (*chat_rpc.UpdateConversationRes, error) {
	var userConvo chat_models.ChatUserConversation
	err := l.svcCtx.DB.Where("conversation_id = ? AND user_id = ?", in.ConversationId, in.UserId).First(&userConvo).Error
	if err != nil {
		// 如果记录不存在，创建新记录
		if err := l.svcCtx.DB.Create(&chat_models.ChatUserConversation{
			UserID:         in.UserId,
			ConversationID: in.ConversationId,
			IsPinned:       in.IsPinned,
			IsHidden:       in.IsDeleted, // 兼容旧的IsDeleted参数
			IsMuted:        false,
			UserReadSeq:    0,
			Version:        1, // 初始版本
		}).Error; err != nil {
			return nil, err
		}
	} else {
		// 如果记录存在，更新记录
		updates := map[string]interface{}{
			"is_hidden": in.IsDeleted, // 兼容旧的IsDeleted参数
		}
		// LastMessage 不再存储在用户会话表中，已移至ChatConversationMeta表
		if in.IsPinned != userConvo.IsPinned {
			updates["is_pinned"] = in.IsPinned
		}

		if err := l.svcCtx.DB.Model(&userConvo).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &chat_rpc.UpdateConversationRes{
		Success: true,
	}, nil
}
