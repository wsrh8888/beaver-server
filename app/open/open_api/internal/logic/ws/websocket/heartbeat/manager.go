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

package heartbeat

import (
	"context"
	"time"

	ws_conn "beaver/app/open/open_api/internal/logic/ws/websocket/conn"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/gorilla/websocket"
)

var logger = beaverlog.New("open_ws_heartbeat")

// HandleClientPing 收到机器人 PING，回复 PONG。
//
// 刻意不调用 coreonline.MarkOnline —— 机器人不是在线用户，不上报用户在线态。
// 这是本服务与用户侧 ws_api 心跳最本质的区别。
func HandleClientPing(client *ws_conn.Client, timestamp int64) {
	if err := client.SafeSendControl(type_struct.WsControlFrame{
		Command:   wsCommandConst.PONG,
		Timestamp: timestamp,
	}); err != nil {
		logger.Error(model.LogMsg{
			Text: "回复PONG失败",
			Data: map[string]any{"timestamp": timestamp, "err": err.Error()},
		})
	}
}

// Manager 服务端主动心跳管理器。
// 维护协议级 WebSocket ping 帧，确保链路不被中间件（如 Nginx/ELB）超时断开。
type Manager struct {
	client     *ws_conn.Client
	robotID    string
	pingPeriod time.Duration
	writeWait  time.Duration
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewManager(client *ws_conn.Client, robotID string, pingPeriod, writeWait time.Duration) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		client:     client,
		robotID:    robotID,
		pingPeriod: pingPeriod,
		writeWait:  writeWait,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (m *Manager) Start() {
	go m.startProtocolPing()
}

func (m *Manager) Stop() {
	m.cancel()
}

// startProtocolPing 协议级 ping（原始 WebSocket 帧）
func (m *Manager) startProtocolPing() {
	ticker := time.NewTicker(m.pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			if err := m.sendProtocolPing(); err != nil {
				logger.Error(model.LogMsg{
					Text: "协议级ping失败",
					Data: map[string]any{"robotId": m.robotID, "err": err.Error()},
				})
				return
			}
		}
	}
}

func (m *Manager) sendProtocolPing() error {
	// 重要：使用 Client 统一的互斥锁，严禁并发写 Conn
	m.client.Mu.Lock()
	defer m.client.Mu.Unlock()

	if err := m.client.Conn.SetWriteDeadline(time.Now().Add(m.writeWait)); err != nil {
		return err
	}
	return m.client.Conn.WriteMessage(websocket.PingMessage, []byte{})
}
