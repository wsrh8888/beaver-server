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

package conn

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	ws_response "beaver/app/ws/ws_api/response"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/gorilla/websocket"
)

// Client 单个 WebSocket 连接封装，带独立写互斥锁，解决并发写问题
type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{Conn: conn}
}

// SafeSend 线程安全发送业务消息（含 content/data 层）
func (c *Client) SafeSend(command wsCommandConst.Command, content type_struct.WsContent) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return ws_response.WsResponse(c.Conn, command, content)
}

// SafeSendControl 线程安全发送控制帧（PING/PONG/ACK，无 content/data 层）
func (c *Client) SafeSendControl(frame type_struct.WsControlFrame) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	data, err := json.Marshal(frame)
	if err != nil {
		fmt.Printf("序列化控制帧失败: %v\n", err)
		return err
	}
	// 设置写入超时，防止受之前写入留下的 deadline 影响
	_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = c.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Printf("发送控制帧失败: %v, 数据: %s\n", err, string(data))
	} else {
		// 仅在调试 PONG 时开启，或者保留少量核心日志
		if frame.Command == wsCommandConst.PONG {
			fmt.Printf("发送 PONG 成功: %s\n", string(data))
		}
	}
	return err
}

// SafeSendJSON 线程安全发送自定义 JSON 数据
func (c *Client) SafeSendJSON(frame type_struct.WsControlFrame) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	data, err := json.Marshal(frame)
	if err != nil {
		fmt.Printf("序列化消息失败: %v\n", err)
		return err
	}
	_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = c.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Printf("发送消息失败: %v\n", err)
	}
	return err
}
