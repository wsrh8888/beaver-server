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

package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	message_models "beaver/app/message/message_models/chat"
	"beaver/app/message/message_rpc/internal/svc"
	"beaver/common/const/mqwsconst"
	"beaver/common/models/ctype"
	"beaver/core/corerocketmq"
	"beaver/remote/agent"
	"beaver/remote/open"
	agenttype "beaver/types/remote/agent"
	opentype "beaver/types/remote/open"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// AgentTaskConsumerLogic 消费 agent_task topic：builtin 数字伙伴消息编排层。
// 从 digital_event_push 投递的 agent_task 信封中取出 payload，校验 Runtime=builtin，
// 取 open 侧绑定的 agentId，拉最近聊天记录拼上下文，最后经 agent_api S2S 调 Python agent
// SendMessage（fire-and-forget）；人设 / 模型 / 工具由 AI 侧按 agentId 解析。
type AgentTaskConsumerLogic struct {
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewAgentTaskConsumerLogic(svcCtx *svc.ServiceContext) *AgentTaskConsumerLogic {
	return &AgentTaskConsumerLogic{
		logger: beaverlog.New("agent_task_consumer", context.Background()),
		svcCtx: svcCtx,
	}
}

// StartConsumer 启动 agent_task 消费者（CLUSTERING，broadcast=false）。
func (l *AgentTaskConsumerLogic) StartConsumer() error {
	mqClient := l.svcCtx.RocketMQ
	if mqClient == nil {
		l.logger.Error(model.LogMsg{Text: "RocketMQ客户端未初始化"})
		return nil
	}

	handler := func(msg *corerocketmq.Message) error {
		// 信封结构：{payload: {messageId, userId, conversationId, content, digitalId}}
		inner, err := extractPayload(msg.Payload)
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "agent_task 信封解析失败", Data: map[string]any{"err": err.Error()}})
			return nil
		}

		digitalID, _ := inner["digitalId"].(string)
		userID, _ := inner["userId"].(string)
		conversationID, _ := inner["conversationId"].(string)
		content, _ := inner["content"].(string)
		messageID, _ := inner["messageId"].(string)
		if digitalID == "" || userID == "" || conversationID == "" {
			l.logger.Error(model.LogMsg{Text: "agent_task 缺少必要字段", Data: map[string]any{"digitalId": digitalID, "userId": userID, "conversationId": conversationID}})
			return nil
		}

		ctx := context.Background()

		// 拉数字伙伴编排配置，校验 Runtime=builtin。
		cfgRes, err := open.GetDigitalConfigApi(ctx, nil, &opentype.GetDigitalConfigReq{DigitalId: digitalID})
		if err != nil || cfgRes == nil {
			l.logger.Error(model.LogMsg{Text: "GetDigitalConfig 失败", Data: map[string]any{"digitalId": digitalID, "err": fmt.Sprintf("%v", err)}})
			return nil
		}
		if cfgRes.Runtime != "builtin" {
			// 非 builtin 模式由其它路径处理，这里直接跳过。
			return nil
		}

		// 编排配置的唯一真相源是 agents 表：这里只透传 agentId，
		// 人设 / 模型 / 工具由 AI 侧按 agentId 取，IM 不再持有并透传 prompt。
		agentID := cfgRes.AgentId
		if agentID == "" {
			l.logger.Error(model.LogMsg{Text: "builtin 数字伙伴未绑定 Agent", Data: map[string]any{"digitalId": digitalID}})
			return nil
		}

		// 拼最近聊天记录上下文。用数字伙伴的 IM 工号识别「自己发的」消息（消息 sender 存的是工号），
		// 不能用 digitalId（那是平台 UUID，与 sender 不同域）。
		historyCtx := l.buildHistoryContext(ctx, conversationID, cfgRes.DigitalUserId)

		// 组装 parts JSON：单个 text part；若有历史则拼到 content 前。
		finalContent := content
		if historyCtx != "" {
			finalContent = historyCtx + "\n\n" + content
		}
		parts := []map[string]string{{"type": "text", "text": finalContent}}
		partsJSON, err := json.Marshal(parts)
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "parts_json 序列化失败", Data: map[string]any{"err": err.Error()}})
			return nil
		}

		// 经 agent_api S2S 调 Python agent SendMessage（fire-and-forget：只取成败，不等终稿）。
		// 信封 code!=0 由 ajax 统一转成 error，这里只看 err。
		_, err = agent.SendMessageApi(ctx, nil, &agenttype.SendMessageReq{
			UserId:         userID,
			ConversationId: conversationID,
			ModelId:        "",
			PartsJson:      string(partsJSON),
			AgentId:        agentID,
			// DigitalId 在 Python 侧是「IM 工号」语义（senderId 兜底 + agent_done 路由回投），
			// 终稿落库拿它当 SendMsg.UserId，所以这里必须传工号而非平台 UUID。
			DigitalId:   cfgRes.DigitalUserId,
			ClientMsgId: messageID,
		})
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "AgentRpc.SendMessage 调用失败", Data: map[string]any{"digitalId": digitalID, "messageId": messageID, "err": err.Error()}})
			return nil
		}

		l.logger.Info(model.LogMsg{Text: "agent_task 投递 agent 成功", Data: map[string]any{"digitalId": digitalID, "messageId": messageID}})
		return nil
	}

	err := mqClient.RegisterConsumer(
		mqwsconst.MqGroupAgentTask,
		l.svcCtx.Config.RocketMQ.Addr,
		mqwsconst.MqTopicAgentTask,
		false, // CLUSTERING：单实例消费
		handler,
	)
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "启动 agent_task Consumer 失败", Data: map[string]any{"err": err.Error()}})
		return err
	}

	l.logger.Info(model.LogMsg{Text: "agent_task Consumer 启动成功"})
	return nil
}

// buildHistoryContext 拉最近 20 条消息（升序）拼成 "角色: 文本" 上下文。
// digitalUserID 是数字伙伴的 IM 工号（user_id 列），用于把「自己发的」消息标成 assistant。
func (l *AgentTaskConsumerLogic) buildHistoryContext(ctx context.Context, conversationID, digitalUserID string) string {
	var msgs []message_models.ChatMessage
	if err := l.svcCtx.DB.Where("conversation_id = ?", conversationID).
		Order("id DESC").Limit(20).Find(&msgs).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询历史消息失败", Data: map[string]any{"conversationId": conversationID, "err": err.Error()}})
		return ""
	}
	if len(msgs) == 0 {
		return ""
	}

	// DB 按 id DESC 取出，反转成升序拼接。
	var b strings.Builder
	for i := len(msgs) - 1; i >= 0; i-- {
		m := msgs[i]
		role := "user"
		if m.SendUserID != nil && *m.SendUserID == digitalUserID {
			role = "assistant"
		}
		text := extractHistoryText(m.Msg)
		if text == "" {
			continue
		}
		fmt.Fprintf(&b, "%s: %s\n", role, text)
	}
	return strings.TrimRight(b.String(), "\n")
}

// extractHistoryText 从 ctype.Msg 提取可读文本（文本/Markdown）。
func extractHistoryText(msg *ctype.Msg) string {
	if msg == nil {
		return ""
	}
	switch msg.Type {
	case ctype.TextMsgType:
		if msg.TextMsg != nil {
			return msg.TextMsg.Content
		}
	case ctype.MarkdownMsgType:
		if msg.MarkdownMsg != nil {
			return msg.MarkdownMsg.Content
		}
	}
	return ""
}

// extractPayload 从信封里取出业务字段。
// Go 投递是 {payload:{...}}；Python 终稿若少包一层，字段就直接在这一层。两种都认。
func extractPayload(envelope map[string]interface{}) (map[string]interface{}, error) {
	if envelope == nil {
		return nil, fmt.Errorf("信封为空")
	}
	v, ok := envelope["payload"]
	if ok && v != nil {
		inner, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("payload 不是 map[string]interface{}")
		}
		return inner, nil
	}
	if _, has := envelope["digitalId"]; has {
		return envelope, nil
	}
	return nil, fmt.Errorf("信封缺少 payload 字段")
}
