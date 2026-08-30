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

type BatchUpdateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUpdateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateConversationLogic {
	return &BatchUpdateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchUpdateConversationLogic) BatchUpdateConversation(in *chat_rpc.BatchUpdateConversationReq) (*chat_rpc.BatchUpdateConversationRes, error) {
	// 开启事务
	tx := l.svcCtx.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// 确保在函数返回时处理事务
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 批量更新或创建会话记录
	for _, userID := range in.UserIds {
		var userConvo chat_models.ChatUserConversation
		err := tx.Where("conversation_id = ? AND user_id = ?", in.ConversationId, userID).First(&userConvo).Error
		if err != nil {
			// 如果记录不存在，创建新记录
			if err := tx.Create(&chat_models.ChatUserConversation{
				UserID:         userID,
				ConversationID: in.ConversationId,
				IsHidden:       false,
				IsPinned:       false,
				IsMuted:        false,
				UserReadSeq:    0,
				Version:        1, // 初始版本
			}).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 如果记录存在，不需要更新LastMessage（已在ChatConversationMeta中）
			// 这里只确保会话没有被隐藏
			if err := tx.Model(&userConvo).Update("is_hidden", false).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &chat_rpc.BatchUpdateConversationRes{
		Success: true,
	}, nil
}
