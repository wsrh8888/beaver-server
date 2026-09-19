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

package ws

import (
	"context"
	"encoding/json"

	ws_conn "beaver/app/open/open_api/internal/logic/ws/websocket/conn"
	"beaver/app/open/open_api/internal/logic/ws/websocket/handler/chat_message"
	"beaver/app/open/open_api/internal/logic/ws/websocket/heartbeat"
	"beaver/app/open/open_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/gorilla/websocket"
)

// inboundFrame 机器人上行帧。
//
// 同时容纳两种形状：
//   - 控制帧：{ "command": "PING", "timestamp": 0 }        —— timestamp 在顶层
//   - 业务帧：{ "command": "CHAT_MESSAGE", "content": {} } —— 对齐 ws_api 的 WsMessage
type inboundFrame struct {
	Command   wsCommandConst.Command `json:"command"`
	Timestamp int64                  `json:"timestamp"`
	Content   *type_struct.WsContent `json:"content"`
}

// HandleRobotMessages 机器人连接的消息循环。
//
// 与用户侧最大的区别：机器人的「回复」不走 WS，统一走 OAPI
// （POST /api/open/robot/v1/send_message），因此上行只有两类帧：
// 心跳，以及流式增量（chat_message_stream_send）。
func HandleRobotMessages(ctx context.Context, svcCtx *svc.ServiceContext, robotID string, client *ws_conn.Client) {
	logger := beaverlog.New("open_ws_handle", ctx)

	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error(model.LogMsg{
					Text: "机器人连接异常关闭",
					Data: map[string]any{"robotId": robotID, "err": err.Error()},
				})
			} else {
				logger.Info(model.LogMsg{
					Text: "机器人连接正常关闭",
					Data: map[string]any{"robotId": robotID},
				})
			}
			break
		}

		var frame inboundFrame
		if err = json.Unmarshal(p, &frame); err != nil {
			logger.Error(model.LogMsg{
				Text: "机器人WS消息解析错误",
				Data: map[string]any{"robotId": robotID, "err": err.Error()},
			})
			continue
		}

		switch frame.Command {
		case wsCommandConst.PING:
			heartbeat.HandleClientPing(client, frame.Timestamp)
		case wsCommandConst.PONG:
			// 服务端协议级 ping 的回执，无需处理
		case wsCommandConst.CHAT_MESSAGE:
			if frame.Content == nil {
				logger.Error(model.LogMsg{
					Text: "机器人业务帧缺少content",
					Data: map[string]any{"robotId": robotID},
				})
				continue
			}
			if err = chat_message.Handle(ctx, svcCtx, robotID, *frame.Content); err != nil {
				logger.Error(model.LogMsg{
					Text: "机器人业务帧处理失败",
					Data: map[string]any{"robotId": robotID, "type": frame.Content.Data.Type, "err": err.Error()},
				})
			}
		default:
			// 不再静默丢弃：带上具体 command，否则龙虾侧调试时会以为帧发出去了
			logger.Error(model.LogMsg{
				Text: "机器人发来未知命令已忽略",
				Data: map[string]any{"robotId": robotID, "command": frame.Command},
			})
		}
	}
}
