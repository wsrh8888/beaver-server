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

package openevent

import (
	"context"
	"errors"

	"beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/core/corerocketmq"
)

// Push 把机器人事件投递到机器人长连接专用 Topic（`ws_bot_push_topic`）。
//
// 所有事件生产侧（chat_rpc / friend_api / open_api 的机器人接口）都走这里，
// 保证 MQ payload 形状一致 —— open_api 侧的长连接消费者依赖这个形状解析并下发。
//
// 参数：
//   - robotID        机器人 IM 用户 ID。既是路由键（open_api 的连接注册表按 robotId 索引），
//                    也是事件归属的机器人。
//   - eventType      事件类型，取值见本包 Event* 常量。会作为 WS 帧的 data.type 下发。
//   - conversationID 会话 ID。消息类事件必填；关系类事件（关注/取关/进群/出群）可传空串。
//   - body           事件体，原样进入 WS 帧的 content.data.body。
//
// 投递后由 open_api 的消费者在本地连接表中查找该机器人：有连接则推送，无连接则丢弃
// （长连接模式下「没连上」即「对方离线」，离线期间事件不补推）。
func Push(ctx context.Context, mq *corerocketmq.Client, robotID, eventType, conversationID string, body map[string]interface{}) error {
	if mq == nil {
		return errors.New("RocketMQ 未初始化")
	}
	if robotID == "" {
		return errors.New("缺少 robotId")
	}
	if eventType == "" {
		return errors.New("缺少 eventType")
	}

	return mq.SendMessage(ctx, mqwsconst.MqTopicWsBot, map[string]interface{}{
		"targetId":       robotID,
		"command":        string(wsCommandConst.CHAT_MESSAGE),
		"type":           eventType,
		"conversationId": conversationID,
		"body":           body,
	})
}
