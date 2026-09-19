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

package chat_message

import (
	"context"
	"errors"

	"beaver/app/chat/chat_models"
	"beaver/app/open/open_api/internal/svc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/conversation"
)

var logger = beaverlog.New("open_ws_chat_message")

// resolveRecipients 解析会话内的收件人，**同时校验机器人确实在该会话中**。
//
// 这一步的校验是必须的：机器人的上行帧里自带 conversationId，
// 若不校验归属，任何已接入的机器人都能给任意用户刷屏。
//
// 收件人算法与 ws_api 的 typing 收件人（getTypingPeerIDs）保持一致：
//   - 私聊：会话 ID 里就是两个用户 ID，取不是机器人的那个
//   - 群聊 / 圈子：以 chat_user_conversations 为准（会话成员表）
func resolveRecipients(ctx context.Context, svcCtx *svc.ServiceContext, robotID, conversationID string) ([]string, error) {
	if conversationID == "" {
		return nil, errors.New("conversationId 不能为空")
	}

	convType, userIDs := conversation.ParseConversationWithType(conversationID)
	if convType == 1 {
		recipients := make([]string, 0, 1)
		inConversation := false
		for _, uid := range userIDs {
			if uid == robotID {
				inConversation = true
				continue
			}
			recipients = append(recipients, uid)
		}
		if !inConversation {
			return nil, errors.New("机器人不在该会话中")
		}
		if len(recipients) == 0 {
			return nil, errors.New("会话内没有其他成员")
		}
		return recipients, nil
	}

	return resolveGroupRecipients(ctx, svcCtx, robotID, conversationID)
}

// resolveGroupRecipients 群聊 / 圈子的收件人解析。
//
// 直接读 chat_user_conversations，而不是调 GroupRpc：
//  1. 一张表同时覆盖群聊与圈子，且能顺带完成「机器人在不在这个会话里」的校验
//  2. 不引入新的 RPC 依赖，也不会有 RPC 失败导致推送整体不可用
func resolveGroupRecipients(ctx context.Context, svcCtx *svc.ServiceContext, robotID, conversationID string) ([]string, error) {
	var rows []chat_models.ChatUserConversation
	if err := svcCtx.DB.WithContext(ctx).Where("conversation_id = ?", conversationID).Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("会话不存在或没有成员")
	}

	recipients := make([]string, 0, len(rows))
	inConversation := false
	for _, row := range rows {
		if row.UserID == robotID {
			inConversation = true
			continue
		}
		recipients = append(recipients, row.UserID)
	}
	if !inConversation {
		return nil, errors.New("机器人不在该会话中")
	}
	return recipients, nil
}
