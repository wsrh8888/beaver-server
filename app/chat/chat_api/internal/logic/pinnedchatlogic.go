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
	"fmt"

	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/app/chat/chat_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type PinnedChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewPinnedChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PinnedChatLogic {
	return &PinnedChatLogic{
		ctx:    ctx,
		logger: beaverlog.New("pinned_chat", ctx),
		svcCtx: svcCtx,
	}
}

func (l *PinnedChatLogic) PinnedChat(req *types.PinnedChatReq) (resp *types.PinnedChatRes, err error) {
	resp = &types.PinnedChatRes{}

	version := l.svcCtx.VersionGen.GetNextVersion("chat_user_conversations", "user_id", req.UserID)

	err = l.svcCtx.DB.Model(&chat_models.ChatUserConversation{}).
		Where("user_id = ? AND conversation_id = ?", req.UserID, req.ConversationID).
		Updates(map[string]interface{}{
			"is_pinned": req.IsPinned,
			"version":   version,
		}).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "更新置顶状态失败",
			Data: map[string]any{
				"userId":         req.UserID,
				"conversationId": req.ConversationID,
				"isPinned":       req.IsPinned,
				"err":            err.Error(),
			},
		})
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "置顶状态更新成功",
		Data: map[string]any{
			"userId":         req.UserID,
			"conversationId": req.ConversationID,
			"isPinned":       req.IsPinned,
			"version":        version,
		},
	})

	go func() {
		l.notifyPinnedUpdate(req.ConversationID, req.UserID, version)
	}()

	return resp, nil
}

func (l *PinnedChatLogic) notifyPinnedUpdate(conversationId, userId string, version int64) {
	defer func() {
		if r := recover(); r != nil {
			l.logger.Error(model.LogMsg{
				Text: "推送置顶通知异常",
				Data: map[string]any{
					"userId":         userId,
					"conversationId": conversationId,
					"panic":          fmt.Sprint(r),
				},
			})
		}
	}()

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

	tableUpdates := []map[string]interface{}{userConversationsUpdate}
	payload := map[string]interface{}{
		"command":  wsCommandConst.CHAT_MESSAGE,
		"type":     wsTypeConst.ChatUserConversationReceive,
		"senderId": userId,
		"targetId": userId,
		"body": map[string]interface{}{
			"tableUpdates": tableUpdates,
		},
		"conversationId": conversationId,
	}
	if err := l.svcCtx.RocketMQ.SendMessage(context.Background(), mqwsconst.MqTopicWs, payload); err != nil {
		l.logger.Error(model.LogMsg{
			Text: "推送置顶通知失败",
			Data: map[string]any{
				"userId":         userId,
				"conversationId": conversationId,
				"err":            err.Error(),
			},
		})
		return
	}

	l.logger.Info(model.LogMsg{
		Text: "推送置顶状态更新通知",
		Data: map[string]any{
			"userId":         userId,
			"conversationId": conversationId,
			"version":        version,
		},
	})
}
