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

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

	"github.com/zeromicro/go-zero/core/logx"
)

type PinnedChatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPinnedChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinnedChatLogic {
	return &PinnedChatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PinnedChatLogic) PinnedChat(req *types.PinnedChatReq) (resp *types.PinnedChatRes, err error) {
	resp = &types.PinnedChatRes{}

	// 获取下一个版本号
	version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)

	// 更新会话置顶状态和版本号
	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{
			"is_pinned": req.IsPinned,
			"version":   version,
		}).Error
	if err != nil {
		l.Logger.Errorf("pinned chat update failed: %v", err)
		return nil, err
	}

	// 发送WS通知给自己（更新本地数据）
	go func() {
		l.notifyPinnedUpdate(req.ConversationID, req.UserID, version)
	}()

	return resp, nil
}

// 发送置顶状态更新通知
func (l *PinnedChatLogic) notifyPinnedUpdate(conversationId, userId string, version int64) {
	defer func() {
		if r := recover(); r != nil {
			l.Logger.Errorf("发送置顶通知时发生panic: %v", r)
		}
	}()

	// 构建用户会话表更新数据
	userConversationsUpdate := map[string]interface{}{
		"table":          "user_conversations",
		"userId":         userId,
		"conversationId": conversationId,
		"data": []map[string]interface{}{
			{
				"version": int32(version),
			},
		},
	}

	// 发送给自己
	tableUpdates := []map[string]interface{}{userConversationsUpdate}
	messageType := wsTypeConst.ChatUserConversationReceive
	payload := map[string]interface{}{
		"command":  wsCommandConst.CHAT_MESSAGE,
		"type":     messageType,
		"senderId": userId,
		"targetId": userId,
		"body": map[string]interface{}{
			"tableUpdates": tableUpdates,
		},
		"conversationId": conversationId,
	}
	l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload)

	l.Logger.Infof("发送置顶状态更新通知: user=%s, conversation=%s, version=%d",
		userId, conversationId, version)
}
