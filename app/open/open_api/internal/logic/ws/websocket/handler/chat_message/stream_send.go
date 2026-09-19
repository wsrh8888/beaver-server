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
	"errors"
	"fmt"

	"beaver/app/open/open_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/beaverlog/model"
)

// StreamBody 机器人上行的流式增量帧体。
//
// 客户端（龙虾）发来的形状：
//
//	{ "command": "CHAT_MESSAGE",
//	  "content": { "data": { "type": "chat_message_stream_send",
//	                         "conversationId": "private_a_b",
//	                         "body": { "streamId": "st_xxx", "delta": "你",
//	                                   "seq": 3, "done": false } } } }
type StreamBody struct {
	StreamID string `json:"streamId"` // 同一轮流式回复的标识，由龙虾生成
	Delta    string `json:"delta"`    // 增量片段（不是全量）
	Seq      int64  `json:"seq"`      // 帧序号，客户端据此丢弃乱序帧
	Done     bool   `json:"done"`     // 是否为本轮最后一帧
}

// HandleStreamSend 把机器人的流式增量转发给会话内的其他成员。
//
// 这是「旁路流式」的服务端实现，关键取舍：
//  1. **不落库、不占 seq** —— 增量是传输态而非存储态。历史记录 / 搜索 / 会话预览
//     读到的永远只有终稿（龙虾走 OAPI send_message 落库），不会读到半截消息。
//  2. **可丢帧** —— 下一帧会补齐，终稿兜底，因此不需要 ACK 与重传。
//  3. **走 ws_push_topic（用户侧通道）** —— 因为收件人是普通 IM 用户，
//     机器人专用 topic（ws_bot_push_topic）只用于「IM → 机器人」这一个方向。
func HandleStreamSend(ctx context.Context, svcCtx *svc.ServiceContext, robotID string, content type_struct.WsContent) error {
	conversationID := content.Data.ConversationID
	if conversationID == "" {
		return errors.New("流式帧缺少 conversationId")
	}
	if len(content.Data.Body) == 0 {
		return errors.New("流式帧缺少 body")
	}

	var body StreamBody
	if err := json.Unmarshal(content.Data.Body, &body); err != nil {
		return fmt.Errorf("流式帧格式错误: %w", err)
	}
	if body.StreamID == "" {
		return errors.New("流式帧缺少 streamId")
	}

	// 解析收件人的同时校验「机器人确实在这个会话里」，防越权刷屏
	recipients, err := resolveRecipients(ctx, svcCtx, robotID, conversationID)
	if err != nil {
		logger.Error(model.LogMsg{
			Text: "流式帧收件人解析失败",
			Data: map[string]any{"robotId": robotID, "conversationId": conversationID, "err": err.Error()},
		})
		return err
	}

	if svcCtx.RocketMQ == nil {
		logger.Error(model.LogMsg{Text: "RocketMQ客户端未初始化，流式帧无法下发"})
		return errors.New("RocketMQ 未初始化")
	}

	outBody := map[string]interface{}{
		"streamId": body.StreamID,
		"senderId": robotID,
		"delta":    body.Delta,
		"seq":      body.Seq,
		"done":     body.Done,
	}

	// 与 sendmsglogic / typing_send 完全相同的信封形状，ws_api 消费者无需改动
	for _, recipientID := range recipients {
		payload := map[string]interface{}{
			"command":        wsCommandConst.CHAT_MESSAGE,
			"type":           wsTypeConst.ChatMessageStreamReceive,
			"senderId":       robotID,
			"targetId":       recipientID,
			"body":           outBody,
			"conversationId": conversationID,
		}
		if err := svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logger.Error(model.LogMsg{
				Text: "流式帧投递失败",
				Data: map[string]any{
					"robotId": robotID, "targetId": recipientID,
					"conversationId": conversationID, "streamId": body.StreamID,
					"err": err.Error(),
				},
			})
		}
	}

	return nil
}
