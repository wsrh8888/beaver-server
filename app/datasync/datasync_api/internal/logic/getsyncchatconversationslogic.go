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
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncChatConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的会话元信息版本
func NewGetSyncChatConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncChatConversationsLogic {
	return &GetSyncChatConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncChatConversationsLogic) GetSyncChatConversations(req *types.GetSyncChatConversationsReq) (resp *types.GetSyncChatConversationsRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 获取用户参与的会话列表
	conversationsResp, err := l.svcCtx.ChatRpc.GetUserConversations(l.ctx, &chat_rpc.GetUserConversationsReq{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("获取用户会话列表失败: %v", err)
		return nil, err
	}

	conversationIDs := make([]string, 0, len(conversationsResp.Conversations))
	for _, conv := range conversationsResp.Conversations {
		conversationIDs = append(conversationIDs, conv.ConversationId)
	}

	if len(conversationIDs) == 0 {
		return &types.GetSyncChatConversationsRes{
			ConversationVersions: []types.ChatConversationVersionItem{},
			ServerTimestamp:      time.Now().UnixMilli(),
		}, nil
	}

	// 获取变更的会话元信息版本
	serverTimestamp := time.Now().UnixMilli()

	convResp, err := l.svcCtx.ChatRpc.GetConversationsListByIds(l.ctx, &chat_rpc.GetConversationsListByIdsReq{
		ConversationIds: conversationIDs,
		Since:           req.Since,
	})
	if err != nil {
		l.Errorf("获取变更的会话版本信息失败: %v", err)
		return nil, err
	}

	// 转换为响应格式，只提取会话元信息版本，确保返回空数组而不是null
	conversationVersions := make([]types.ChatConversationVersionItem, 0)
	if convResp.Conversations != nil {
		for _, conv := range convResp.Conversations {
			conversationVersions = append(conversationVersions, types.ChatConversationVersionItem{
				ConversationID: conv.ConversationId,
				Version:        conv.Version,
			})
		}
	}

	return &types.GetSyncChatConversationsRes{
		ConversationVersions: conversationVersions,
		ServerTimestamp:      serverTimestamp,
	}, nil
}
