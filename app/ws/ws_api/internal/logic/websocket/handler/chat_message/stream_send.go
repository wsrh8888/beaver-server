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
	"encoding/json"
	"fmt"
	"net/http"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/beaverlog/model"
)

// HandleStreamSend 把一段增量转发给该看这场回复的在线设备。
//
// 不落库、不占 seq。丢帧不影响历史，终稿仍走正式消息。
// 收件人是当前用户自己的设备，外加这个会话里的其他人。
// 这样私聊、群、以及没有标准会话的助手入口都能慢慢推。
func HandleStreamSend(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *types.WsReq,
	_ *http.Request,
	_ *ws_conn.Client,
	conversationID string,
	bodyRaw json.RawMessage,
) error {
	var body type_struct.BodyStream
	if err := json.Unmarshal(bodyRaw, &body); err != nil {
		return fmt.Errorf("流式帧格式错误: %w", err)
	}
	if conversationID == "" {
		conversationID = body.ConversationID
	}
	if conversationID == "" {
		return fmt.Errorf("conversationId 不能为空")
	}
	if body.StreamID == "" {
		return fmt.Errorf("streamId 不能为空")
	}
	if svcCtx.RocketMQ == nil {
		return fmt.Errorf("RocketMQ 未初始化")
	}

	recipients := streamRecipients(ctx, svcCtx, req.UserID, conversationID)
	outBody := map[string]interface{}{
		"streamId": body.StreamID,
		"senderId": req.UserID,
		"delta":    body.Delta,
		"seq":      body.Seq,
		"done":     body.Done,
	}

	for _, recipientID := range recipients {
		payload := map[string]interface{}{
			"command":        wsCommandConst.CHAT_MESSAGE,
			"type":           wsTypeConst.ChatMessageStreamReceive,
			"senderId":       req.UserID,
			"targetId":       recipientID,
			"body":           outBody,
			"conversationId": conversationID,
		}
		if err := svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logger.Error(model.LogMsg{
				Text: "推送流式增量失败",
				Data: map[string]any{
					"sender": req.UserID, "target": recipientID,
					"conversationId": conversationID, "streamId": body.StreamID,
					"seq": body.Seq, "err": err.Error(),
				},
			})
		}
	}
	return nil
}

func streamRecipients(ctx context.Context, svcCtx *svc.ServiceContext, currentUserID, conversationID string) []string {
	recipients := []string{currentUserID}
	peers, err := getTypingPeerIDs(ctx, svcCtx, currentUserID, conversationID)
	if err != nil {
		return recipients
	}
	seen := map[string]struct{}{currentUserID: {}}
	for _, peerID := range peers {
		if _, ok := seen[peerID]; ok {
			continue
		}
		seen[peerID] = struct{}{}
		recipients = append(recipients, peerID)
	}
	return recipients
}
