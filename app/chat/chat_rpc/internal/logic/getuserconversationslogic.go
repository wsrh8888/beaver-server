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

type GetUserConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserConversationsLogic {
	return &GetUserConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserConversationsLogic) GetUserConversations(in *chat_rpc.GetUserConversationsReq) (*chat_rpc.GetUserConversationsRes, error) {
	var userConversations []chat_models.ChatUserConversation
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).Find(&userConversations).Error
	if err != nil {
		l.Errorf("查询用户会话失败: %v", err)
		return nil, err
	}

	var conversations []*chat_rpc.ConversationItem
	for _, uc := range userConversations {
		// 查询会话类型
		var conversation chat_models.ChatConversationMeta
		err := l.svcCtx.DB.Where("conversation_id = ?", uc.ConversationID).First(&conversation).Error
		if err != nil {
			l.Errorf("查询会话信息失败: %v", err)
			continue
		}

		conversations = append(conversations, &chat_rpc.ConversationItem{
			ConversationId: uc.ConversationID,
			Type:           int32(conversation.Type),
		})
	}

	return &chat_rpc.GetUserConversationsRes{
		Conversations: conversations,
	}, nil
}
