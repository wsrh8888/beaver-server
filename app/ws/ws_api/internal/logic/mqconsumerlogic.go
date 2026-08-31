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

package logic

import (
	"context"
	"encoding/json"
	"fmt"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/core/corerocketmq"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type MqConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewMqConsumerLogic(svcCtx *svc.ServiceContext) *MqConsumerLogic {
	return &MqConsumerLogic{
		logger: beaverlog.New("mq_consumer", context.Background()),
		svcCtx: svcCtx,
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

// StartConsumer 启动 RocketMQ 消费者
func (l *MqConsumerLogic) StartConsumer() error {
	mqClient := l.svcCtx.RocketMQ
	if mqClient == nil {
		l.logger.Error(model.LogMsg{Text: "RocketMQ客户端未初始化"})
		return nil
	}

	handler := func(msg *corerocketmq.Message) error {
		targetID, err := payloadString(msg.Payload, "targetId")
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "MQ消息格式错误", Data: map[string]any{"field": "targetId", "err": err.Error()}})
			return nil
		}
		command, err := payloadString(msg.Payload, "command")
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "MQ消息格式错误", Data: map[string]any{"field": "command", "err": err.Error()}})
			return nil
		}
		msgType, err := payloadString(msg.Payload, "type")
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "MQ消息格式错误", Data: map[string]any{"field": "type", "err": err.Error()}})
			return nil
		}
		// 通知/好友/群等非聊天类推送允许 conversationId 为空
		conversationID := ""
		if v, ok := msg.Payload["conversationId"]; ok && v != nil {
			if s, ok := v.(string); ok {
				conversationID = s
			}
		}

		body, ok := msg.Payload["body"]
		if !ok || body == nil {
			l.logger.Error(model.LogMsg{Text: "MQ消息缺少body字段", Data: map[string]any{"targetId": targetID}})
			return nil
		}
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "序列化body失败", Data: map[string]any{"targetId": targetID, "err": err.Error()}})
			return err
		}

		content := type_struct.WsContent{
			Data: type_struct.WsData{
				Type:           wsTypeConst.Type(msgType),
				Body:           bodyBytes,
				ConversationID: conversationID,
			},
		}

		ws_conn.SendMsgToUser(targetID, wsCommandConst.Command(command), content)
		return nil
	}

	err := mqClient.RegisterConsumer(
		mqwsconst.MqGroupWs,
		l.svcCtx.Config.RocketMQ.Addr,
		mqwsconst.MqTopicWs,
		true, // 广播模式：每个 WS 实例都消费，才能推送到本机在线连接
		handler,
	)

	if err != nil {
		l.logger.Error(model.LogMsg{Text: "启动RocketMQConsumer失败", Data: map[string]any{"err": err.Error()}})
		return err
	}

	l.logger.Info(model.LogMsg{Text: "WSAPIRocketMQConsumer启动成功"})
	return nil
}
