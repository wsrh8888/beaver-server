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
	"strings"
	"sync"

	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserOnlineWsMap key: userID_deviceType → UserWsInfo
var UserOnlineWsMap = make(map[string]*UserWsInfo)
var WsMapMutex sync.RWMutex

type UserWsInfo struct {
	WsClientMap map[string]*Client // key: conn.RemoteAddr().String()
}

func GetUserKey(userID, deviceType string) string {
	return userID + "_" + deviceType
}

func IsGroupChat(conversationID string) bool {
	return !strings.Contains(conversationID, "_")
}

func GetRecipientIdFromConversationID(conversationID, userID string) string {
	ids := strings.Split(conversationID, "_")
	if len(ids) != 2 {
		return ""
	}
	if ids[0] == userID {
		return ids[1]
	}
	return ids[0]
}

// SendMsgToUser 发送消息给指定用户的所有在线设备
// 先在读锁下收集目标连接，再释放读锁后发送，避免在锁内写 WebSocket
func SendMsgToUser(targetUserID string, command wsCommandConst.Command, content type_struct.WsContent) {
	WsMapMutex.RLock()
	var targets []*Client
	for userKey, userInfo := range UserOnlineWsMap {
		if strings.HasPrefix(userKey, targetUserID+"_") {
			for _, client := range userInfo.WsClientMap {
				targets = append(targets, client)
			}
		}
	}
	WsMapMutex.RUnlock()

	for _, client := range targets {
		if err := client.SafeSend(command, content); err != nil {
			logx.Errorf("发送消息给用户 %s 失败: %v", targetUserID, err)
		}
	}
}

// SendMsgToReceiverAndSyncToSender 发消息给接收者，并同步给发送者的其他设备
func SendMsgToReceiverAndSyncToSender(
	revUserID string,
	sendUserID string,
	command wsCommandConst.Command,
	receiverContent type_struct.WsContent,
	senderSyncContent type_struct.WsContent,
	excludeAddr string,
) {
	SendMsgToUser(revUserID, command, receiverContent)

	WsMapMutex.RLock()
	var senderTargets []*Client
	for userKey, userInfo := range UserOnlineWsMap {
		if strings.HasPrefix(userKey, sendUserID+"_") {
			for addr, client := range userInfo.WsClientMap {
				if addr != excludeAddr {
					senderTargets = append(senderTargets, client)
				}
			}
		}
	}
	WsMapMutex.RUnlock()

	for _, client := range senderTargets {
		if err := client.SafeSend(command, senderSyncContent); err != nil {
			logx.Errorf("同步消息给发送者 %s 失败: %v", sendUserID, err)
		}
	}
}
