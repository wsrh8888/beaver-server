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

type GetConversationsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsListByIdsLogic {
	return &GetConversationsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationsListByIdsLogic) GetConversationsListByIds(in *chat_rpc.GetConversationsListByIdsReq) (*chat_rpc.GetConversationsListByIdsRes, error) {
	// 构建查询条件
	query := l.svcCtx.DB.Where("conversation_id IN (?)", in.ConversationIds)

	// 如果提供了since参数，只返回版本有变更的记录
	if in.Since > 0 {
		query = query.Where("version >= ?", in.Since)
	}

	// 查询指定会话的完整信息
	var conversations []chat_models.ChatConversationMeta
	err := query.Find(&conversations).Error
	if err != nil {
		l.Errorf("查询会话信息失败: %v", err)
		return nil, err
	}

	var conversationList []*chat_rpc.ConversationListById
	for _, conv := range conversations {
		conversationList = append(conversationList, &chat_rpc.ConversationListById{
			ConversationId: conv.ConversationID,
			Type:           int32(conv.Type),
			Seq:            conv.MaxSeq,
			Version:        conv.Version,
		})
	}

	return &chat_rpc.GetConversationsListByIdsRes{
		Conversations: conversationList,
	}, nil
}
