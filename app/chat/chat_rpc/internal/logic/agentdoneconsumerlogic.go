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

	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/common/const/mqwsconst"
	"beaver/common/models/ctype"
	"beaver/core/corerocketmq"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/google/uuid"
)

// AgentDoneConsumerLogic 消费 agent_done topic：agent 跑完 LLM 后把终稿回投给 chat 落库。
// payload 结构：{digitalId, conversationId, content, userId}。
// 调 ChatRpc.SendMsg 以 digitalId 为发送者、DeviceId=open_digital（跳过好友校验）落一条文本消息。
type AgentDoneConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAgentDoneConsumerLogic(svcCtx *svc.ServiceContext) *AgentDoneConsumerLogic {
	return &AgentDoneConsumerLogic{
		logger: beaverlog.New("agent_done_consumer", context.Background()),
		svcCtx: svcCtx,
	}
}

// StartConsumer 启动 agent_done 消费者（CLUSTERING，broadcast=false）。
func (l *AgentDoneConsumerLogic) StartConsumer() error {
	mqClient := l.svcCtx.RocketMQ
	if mqClient == nil {
		l.logger.Error(model.LogMsg{Text: "RocketMQ客户端未初始化"})
		return nil
	}

	handler := func(msg *corerocketmq.Message) error {
		// 信封结构：{payload: {digitalId, conversationId, content, userId}}
		inner, err := extractPayload(msg.Payload)
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "agent_done 信封解析失败", Data: map[string]any{"err": err.Error()}})
			return nil
		}

		digitalID, _ := inner["digitalId"].(string)
		conversationID, _ := inner["conversationId"].(string)
		content, _ := inner["content"].(string)
		if digitalID == "" || conversationID == "" || content == "" {
			l.logger.Error(model.LogMsg{Text: "agent_done 缺少必要字段", Data: map[string]any{"digitalId": digitalID, "conversationId": conversationID}})
			return nil
		}

		ctx := context.Background()

		// 以数字伙伴身份落一条文本消息；DeviceId=open_digital 让 sendmsglogic 跳过好友校验。
		sendLogic := NewSendMsgLogic(ctx, l.svcCtx)
		_, err = sendLogic.SendMsg(&chat_rpc.SendMsgReq{
			UserId:         digitalID,
			ConversationId: conversationID,
			MessageId:      uuid.New().String(),
			DeviceId:       "open_digital",
			Msg: &chat_rpc.Msg{
				Type:    uint32(ctype.TextMsgType),
				TextMsg: &chat_rpc.TextMsg{Content: content},
			},
		})
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "agent_done 落库失败", Data: map[string]any{"digitalId": digitalID, "conversationId": conversationID, "err": err.Error()}})
			// 不阻塞队列：返回 nil。
			return nil
		}

		l.logger.Info(model.LogMsg{Text: "agent_done 落库成功", Data: map[string]any{"digitalId": digitalID, "conversationId": conversationID}})
		return nil
	}

	err := mqClient.RegisterConsumer(
		mqwsconst.MqGroupAgentDone,
		l.svcCtx.Config.RocketMQ.Addr,
		mqwsconst.MqTopicAgentDone,
		false, // CLUSTERING：单实例消费
		handler,
	)
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "启动 agent_done Consumer 失败", Data: map[string]any{"err": err.Error()}})
		return err
	}

	l.logger.Info(model.LogMsg{Text: "agent_done Consumer 启动成功"})
	return nil
}
