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

type GetSyncChatUserConversationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的用户会话设置版本
func NewGetSyncChatUserConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncChatUserConversationsLogic {
	return &GetSyncChatUserConversationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncChatUserConversationsLogic) GetSyncChatUserConversations(req *types.GetSyncChatUserConversationsReq) (resp *types.GetSyncChatUserConversationsRes, err error) {
	userId := req.UserID
	if userId == "" {
		l.Errorf("用户ID为空")
		return nil, errors.New("用户ID不能为空")
	}

	// 直接获取用户的所有会话设置版本（一个RPC调用搞定）
	serverTimestamp := time.Now().UnixMilli()

	userConvResp, err := l.svcCtx.ChatRpc.GetUserConversationVersions(l.ctx, &chat_rpc.GetUserConversationVersionsReq{
		UserId: userId,
		Since:  req.Since,
	})
	if err != nil {
		l.Errorf("获取用户会话设置版本失败: %v", err)
		return nil, err
	}

	// 转换为响应格式，确保返回空数组而不是null
	userConversationVersions := make([]types.ChatUserConversationVersionItem, 0)
	if userConvResp.UserConversationVersions != nil {
		for _, userConv := range userConvResp.UserConversationVersions {
			userConversationVersions = append(userConversationVersions, types.ChatUserConversationVersionItem{
				ConversationID: userConv.ConversationId,
				Version:        userConv.Version,
			})
		}
	}

	return &types.GetSyncChatUserConversationsRes{
		UserConversationVersions: userConversationVersions,
		ServerTimestamp:          serverTimestamp,
	}, nil
}
