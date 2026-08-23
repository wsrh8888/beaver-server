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
	"strconv"

	"beaver/app/chat/chat_models"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	"beaver/core/coreonline"
	"beaver/core/corepush"

	"github.com/zeromicro/go-zero/core/logx"
)

func (l *SendMsgLogic) sendOfflinePushIfNeeded(
	conversationID, senderID string,
	chatModel chat_models.ChatMessage,
	recipientIDs []string,
) {
	if l.svcCtx.PushSender == nil || !l.svcCtx.PushSender.Enabled() {
		return
	}

	sender, err := l.getSenderInfo(chatModel)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("离线推送获取发送者失败: %v", err)
		return
	}

	title := sender.NickName
	if title == "" {
		title = "新消息"
	}
	body := chatModel.MsgPreview
	if body == "" {
		body = "你收到一条新消息"
	}

	data := map[string]string{
		"type":           "chat_message",
		"conversationId": conversationID,
		"messageId":      chatModel.MessageID,
		"seq":            strconv.FormatInt(chatModel.Seq, 10),
		"senderId":       senderID,
	}

	for _, recipientID := range recipientIDs {
		if recipientID == senderID {
			continue
		}
		if coreonline.IsOnline(l.svcCtx.Redis, recipientID) {
			continue
		}

		var userConvo chat_models.ChatUserConversation
		if err := l.svcCtx.DB.Where("user_id = ? AND conversation_id = ?", recipientID, conversationID).
			First(&userConvo).Error; err == nil && userConvo.IsMuted {
			continue
		}

		res, err := l.svcCtx.NotificationRpc.ListPushTokens(context.Background(), &notification_rpc.ListPushTokensReq{
			UserId: recipientID,
		})
		if err != nil || len(res.Tokens) == 0 {
			continue
		}

		tokens := make([]corepush.PushToken, 0, len(res.Tokens))
		for _, t := range res.Tokens {
			tokens = append(tokens, corepush.PushToken{
				DeviceID:     t.DeviceId,
				PushToken:    t.PushToken,
				PushPlatform: t.PushPlatform,
			})
		}

		l.svcCtx.PushSender.SendToTokens(context.Background(), tokens, corepush.Message{
			Title: title,
			Body:  body,
			Data:  data,
		})
	}
}
