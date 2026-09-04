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
	"net/http"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/logic/websocket/handler/chat_message"
	"beaver/app/ws/ws_api/internal/logic/websocket/heartbeat"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/gorilla/websocket"
)

var logger = beaverlog.New("ws_handle")

func HandleWebSocketMessages(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, deviceGroup string, r *http.Request, client *ws_conn.Client) {
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error(model.LogMsg{
					Text: "WebSocket连接异常关闭",
					Data: map[string]interface{}{"userId": req.UserID, "err": err.Error()},
				})
			} else {
				logger.Info(model.LogMsg{
					Text: "WebSocket连接正常关闭",
					Data: map[string]interface{}{"userId": req.UserID},
				})
			}
			break
		}

		var wsMessage type_struct.WsMessage
		if err = json.Unmarshal(p, &wsMessage); err != nil {
			logger.Error(model.LogMsg{
				Text: "WS消息解析错误",
				Data: map[string]interface{}{"userId": req.UserID, "err": err.Error()},
			})
			continue
		}

		source := map[string]string{
			"platform":   req.Platform,
			"deviceId":   req.DeviceID,
			"remoteAddr": client.Conn.RemoteAddr().String(),
			"userAgent":  r.Header.Get("User-Agent"),
		}

		if wsMessage.Command == "" {
			logger.Info(model.LogMsg{
				Text: "收到WS消息",
				Data: map[string]interface{}{
					"userId":  req.UserID,
					"content": json.RawMessage(p),
					"source":  source,
				},
			})
			continue
		}

		cmd := wsCommandConst.Command(wsMessage.Command)
		logger.Info(model.LogMsg{
			Text: "收到WS消息",
			Data: map[string]interface{}{
				"userId": req.UserID,
				"content": map[string]interface{}{
					"command": wsMessage.Command,
					"content": wsMessage.Content,
				},
				"source": source,
			},
		})

		// 控制帧：PING/PONG 直接处理，不发 ACK
		switch cmd {
		case wsCommandConst.PING:
			heartbeat.HandleClientPing(svcCtx.Redis, req.UserID, deviceGroup, svcCtx.InstanceID, client, wsMessage.Content.Timestamp)
			continue
		case wsCommandConst.PONG:
			continue
		case wsCommandConst.USER_PROFILE, wsCommandConst.NOTIFICATION, wsCommandConst.EMOJI:
			// 仅服务端推送，客户端不应发送
			logger.Info(model.LogMsg{
				Text: "客户端不应发送此命令",
				Data: map[string]interface{}{"userId": req.UserID, "command": cmd},
			})
			continue
		}

		// 业务命令：立即发送 ACK（表示服务端已收到），再处理
		msgId := wsMessage.Content.MessageID
		client.SafeSendControl(type_struct.WsControlFrame{
			Command:   wsCommandConst.ACK,
			MessageID: msgId,
		})

		var handlerErr error
		switch cmd {
		case wsCommandConst.CHAT_MESSAGE:
			handlerErr = chat_message.Handle(ctx, svcCtx, req, r, client, wsMessage.Content)
		default:
			logger.Info(model.LogMsg{
				Text: "未支持的命令类型",
				Data: map[string]interface{}{"userId": req.UserID, "command": wsMessage.Command},
			})
		}

		if handlerErr != nil {
			logger.Error(model.LogMsg{
				Text: "处理命令失败",
				Data: map[string]interface{}{"userId": req.UserID, "command": cmd, "err": handlerErr.Error()},
			})
		}
	}
}
