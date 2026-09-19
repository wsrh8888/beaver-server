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
	"sync"

	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

var logger = beaverlog.New("open_ws_conn")

// BotOnlineWsMap 机器人连接注册表，key: robotID（= OpenRobot.RobotID = IM 用户 ID）
//
// 刻意与用户侧 ws_api 的 UserOnlineWsMap 完全分离，原因：
//  1. 机器人不是「在线用户」，混入会污染 coreonline 在线态与后台在线统计
//  2. 用户侧推送是 O(N) 全表扫描，混入只会放大开销
//
// 索引直接用 robotID 而非前缀匹配，命中是 O(1)。
var (
	BotOnlineWsMap = make(map[string]*BotWsInfo)
	BotWsMapMutex  sync.RWMutex
)

// BotWsInfo 单个机器人持有的连接集合
type BotWsInfo struct {
	WsClientMap map[string]*Client // key: conn.RemoteAddr().String()
}

// Register 注册一条机器人连接。
// 同一机器人重复建连时踢掉旧连接（一个机器人实例一条连接），返回被踢掉的数量。
func Register(robotID string, client *Client) int {
	BotWsMapMutex.Lock()
	defer BotWsMapMutex.Unlock()

	addr := client.Conn.RemoteAddr().String()
	kicked := 0

	info, ok := BotOnlineWsMap[robotID]
	if !ok {
		BotOnlineWsMap[robotID] = &BotWsInfo{
			WsClientMap: map[string]*Client{addr: client},
		}
		return 0
	}

	for oldAddr, oldClient := range info.WsClientMap {
		if oldAddr == addr {
			continue
		}
		logger.Info(model.LogMsg{
			Text: "机器人重复建连关闭旧连接",
			Data: map[string]any{"robotId": robotID, "oldAddr": oldAddr},
		})
		_ = oldClient.Conn.Close()
		delete(info.WsClientMap, oldAddr)
		kicked++
	}
	info.WsClientMap[addr] = client
	return kicked
}

// Unregister 注销一条机器人连接；连接全部断开时删除该机器人的条目。
func Unregister(robotID, addr string) {
	BotWsMapMutex.Lock()
	defer BotWsMapMutex.Unlock()

	info, ok := BotOnlineWsMap[robotID]
	if !ok {
		return
	}
	delete(info.WsClientMap, addr)
	if len(info.WsClientMap) == 0 {
		delete(BotOnlineWsMap, robotID)
	}
}

// IsOnline 判断机器人当前是否有活跃连接（供事件生产侧判断投递方式）
func IsOnline(robotID string) bool {
	BotWsMapMutex.RLock()
	defer BotWsMapMutex.RUnlock()

	info, ok := BotOnlineWsMap[robotID]
	return ok && len(info.WsClientMap) > 0
}

// SendMsgToRobot 把事件推送给指定机器人的所有活跃连接。
// 先在读锁下收集目标连接，再释放锁后发送，避免在锁内写 WebSocket。
func SendMsgToRobot(robotID string, command wsCommandConst.Command, content type_struct.WsContent) {
	BotWsMapMutex.RLock()
	var targets []*Client
	if info, ok := BotOnlineWsMap[robotID]; ok {
		for _, client := range info.WsClientMap {
			targets = append(targets, client)
		}
	}
	BotWsMapMutex.RUnlock()

	for _, client := range targets {
		if err := client.SafeSend(command, content); err != nil {
			logger.Error(model.LogMsg{
				Text: "推送事件给机器人失败",
				Data: map[string]any{"robotId": robotID, "err": err.Error()},
			})
		}
	}
}
