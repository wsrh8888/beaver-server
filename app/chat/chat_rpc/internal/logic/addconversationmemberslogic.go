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

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddConversationMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddConversationMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddConversationMembersLogic {
	return &AddConversationMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddConversationMembersLogic) AddConversationMembers(in *chat_rpc.AddConversationMembersReq) (*chat_rpc.AddConversationMembersRes, error) {
	if in.ConversationId == "" {
		return nil, errors.New("会话ID不能为空")
	}

	if len(in.UserIds) == 0 {
		return nil, errors.New("用户列表不能为空")
	}

	// 检查会话是否存在
	var conversation chat_models.ChatConversationMeta
	if err := l.svcCtx.DB.Where("conversation_id = ?", in.ConversationId).First(&conversation).Error; err != nil {
		l.Errorf("会话不存在: conversationId=%s, error=%v", in.ConversationId, err)
		return nil, errors.New("会话不存在")
	}

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

	// 为新成员创建用户会话关系记录
	for _, userId := range in.UserIds {
		// 检查是否已存在
		var existing chat_models.ChatUserConversation
		err := tx.Where("conversation_id = ? AND user_id = ?", in.ConversationId, userId).First(&existing).Error
		if err == nil {
			// 已存在，跳过
			continue
		}

		// 创建新的用户会话关系
		userConversation := chat_models.ChatUserConversation{
			UserID:         userId,
			ConversationID: in.ConversationId,
			IsPinned:       false,
			IsMuted:        false,
			UserReadSeq:    conversation.MaxSeq, // 新成员的已读序列号设为当前最大序列号
			Version:        1,
		}

		if err := tx.Create(&userConversation).Error; err != nil {
			l.Errorf("创建用户会话关系失败: userId=%s, conversationId=%s, error=%v", userId, in.ConversationId, err)
			tx.Rollback()
			return nil, errors.New("创建用户会话关系失败")
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	l.Infof("成功添加会话成员: conversationId=%s, userIds=%v", in.ConversationId, in.UserIds)

	return &chat_rpc.AddConversationMembersRes{
		Success: true,
	}, nil
}
