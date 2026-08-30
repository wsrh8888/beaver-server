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

package ws_response

import (
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	utils "beaver/utils/rand"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Response struct {
	Code       int                    `json:"code"`
	Command    wsCommandConst.Command `json:"command"`
	Content    type_struct.WsContent  `json:"content"`
	MessageID  string                 `json:"messageId"`
	ServerTime int64                  `json:"serverTime"`
}

func WsResponse(conn *websocket.Conn, command wsCommandConst.Command, content type_struct.WsContent) error {
	code := 0

	response := Response{
		Command:    command,
		Code:       code,
		Content:    content,
		MessageID:  utils.GenerateRandomString(8),
		ServerTime: time.Now().Unix(),
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		logx.Errorf("序列化WebSocket响应失败: %v", err)
		return err
	}

	// 设置写入超时
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		logx.Errorf("设置WebSocket写入超时失败: %v", err)
		return err
	}

	if err := conn.WriteMessage(websocket.TextMessage, responseJSON); err != nil {
		logx.Errorf("发送WebSocket消息失败: %v", err)
		return err
	}

	return nil
}
