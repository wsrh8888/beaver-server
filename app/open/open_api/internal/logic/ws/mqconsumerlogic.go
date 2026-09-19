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
	"fmt"
	"time"

	ws_conn "beaver/app/open/open_api/internal/logic/ws/websocket/conn"
	"beaver/app/open/open_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/core/corerocketmq"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// MqConsumerLogic 消费 IM 事件并推送给在线机器人连接。
//
// 事件由 chat_rpc 在消息落库后投递到 MqTopicWsBot，payload 形如：
//
//	{
//	  "targetId":       "<robotId>",              // 路由键：机器人的 IM 用户 ID
//	  "command":        "CHAT_MESSAGE",
//	  "type":           "im.message.receive",     // 事件类型
//	  "conversationId": "private_a_b",
//	  "body":           { ... }                   // 事件体（sender_id / content / ...）
//	}
type MqConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewMqConsumerLogic(svcCtx *svc.ServiceContext) *MqConsumerLogic {
	return &MqConsumerLogic{
		svcCtx: svcCtx,
		logger: beaverlog.New("open_ws_mq_consumer", context.Background()),
	}
}

func payloadString(payload map[string]interface{}, key string) (string, error) {
	v, ok := payload[key]
	if !ok || v == nil {
		return "", fmt.Errorf("缺少字段 %s", key)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", fmt.Errorf("字段 %s 无效", key)
	}
	return s, nil
}

// StartConsumer 启动机器人事件消费者。
//
// 广播模式：每个 open_api 实例都消费，谁持有该机器人的连接谁推送。
func (l *MqConsumerLogic) StartConsumer() error {
	mqClient := l.svcCtx.RocketMQ
	if mqClient == nil {
		l.logger.Error(model.LogMsg{Text: "RocketMQ客户端未初始化"})
		return nil
	}

	handler := func(msg *corerocketmq.Message) error {
		targetID, err := payloadString(msg.Payload, "targetId")
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "机器人事件缺少路由字段", Data: map[string]any{"field": "targetId", "err": err.Error()}})
			return nil
		}

		// 本实例没有该机器人的连接：直接丢弃。
		// 长连接模式下「没连上」即「对方离线」，事件不补推（重连后由机器人自己拉取）。
		if !ws_conn.IsOnline(targetID) {
			l.logger.Info(model.LogMsg{
				Text: "机器人当前无连接事件已丢弃",
				Data: map[string]any{"robotId": targetID},
			})
			return nil
		}

		command, err := payloadString(msg.Payload, "command")
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "机器人事件缺少字段", Data: map[string]any{"field": "command", "err": err.Error()}})
			return nil
		}
		eventType, err := payloadString(msg.Payload, "type")
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "机器人事件缺少字段", Data: map[string]any{"field": "type", "err": err.Error()}})
			return nil
		}

		conversationID := ""
		if v, ok := msg.Payload["conversationId"]; ok && v != nil {
			if s, ok := v.(string); ok {
				conversationID = s
			}
		}

		body, ok := msg.Payload["body"]
		if !ok || body == nil {
			l.logger.Error(model.LogMsg{Text: "机器人事件缺少body", Data: map[string]any{"robotId": targetID}})
			return nil
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "序列化机器人事件body失败", Data: map[string]any{"robotId": targetID, "err": err.Error()}})
			return err
		}

		content := type_struct.WsContent{
			Timestamp: time.Now().UnixMilli(),
			Data: type_struct.WsData{
				// type 承载事件类型（im.message.receive / im.message.receive.group_at）
				Type:           wsTypeConst.Type(eventType),
				ConversationID: conversationID,
				Body:           bodyBytes,
			},
		}

		ws_conn.SendMsgToRobot(targetID, wsCommandConst.Command(command), content)
		return nil
	}

	err := mqClient.RegisterConsumer(
		mqwsconst.MqGroupWsBot,
		l.svcCtx.Config.RocketMQ.Addr,
		mqwsconst.MqTopicWsBot,
		true, // 广播模式：每个实例都消费，才能推送到本机持有的连接
		handler,
	)
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "启动机器人事件消费者失败", Data: map[string]any{"err": err.Error()}})
		return err
	}

	l.logger.Info(model.LogMsg{Text: "机器人事件消费者启动成功"})
	return nil
}
