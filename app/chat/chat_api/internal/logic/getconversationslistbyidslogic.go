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
	"time"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取会话数据
func NewGetConversationsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsListByIdsLogic {
	return &GetConversationsListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetConversationsListByIdsLogic) GetConversationsListByIds(req *types.GetConversationsListByIdsReq) (resp *types.GetConversationsListByIdsRes, err error) {
	// 直接从数据库查询会话完整信息
	var conversations []chat_models.ChatConversationMeta
	err = l.svcCtx.DB.Where("conversation_id IN (?)", req.ConversationIds).Find(&conversations).Error
	if err != nil {
		l.Errorf("查询会话信息失败: %v", err)
		return nil, err
	}

	// 转换数据库模型为API响应
	conversationList := make([]types.ConversationById, 0, len(conversations))
	for _, conv := range conversations {
		conversationList = append(conversationList, types.ConversationById{
			ConversationID: conv.ConversationID,
			Type:           conv.Type,
			MaxSeq:         conv.MaxSeq,
			LastMessage:    conv.LastMessage,
			Version:        conv.Version,
			CreatedAt:      time.Time(conv.CreatedAt).Unix(),
			UpdatedAt:      time.Time(conv.UpdatedAt).Unix(),
		})
	}

	return &types.GetConversationsListByIdsRes{
		Conversations: conversationList,
	}, nil
}
