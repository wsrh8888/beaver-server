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
	"strings"

	chat_models "beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/agent/agent_rpc/agent"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	"beaver/common/const/mqwsconst"
	"beaver/common/models/ctype"
	"beaver/core/corerocketmq"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// AgentTaskConsumerLogic 消费 agent_task topic：builtin 数字伙伴消息编排层。
// 从 digital_event_push 投递的 agent_task 信封中取出 payload，校验 Runtime=builtin，
// 取 BuiltinConfig，拉最近聊天记录拼上下文，最后调用 Python agent_rpc SendMessage（fire-and-forget）。
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
		cfgRes, err := l.svcCtx.OpenRpc.GetDigitalConfig(ctx, &open_rpc.GetDigitalConfigReq{DigitalId: digitalID})
		if err != nil || cfgRes == nil {
			l.logger.Error(model.LogMsg{Text: "GetDigitalConfig 失败", Data: map[string]any{"digitalId": digitalID, "err": fmt.Sprintf("%v", err)}})
			return nil
		}
		if cfgRes.Runtime != "builtin" {
			// 非 builtin 模式由其它路径处理，这里直接跳过。
			return nil
		}

		// 取 AI 应用配置（platform 层，渠道无关通用配置）。
		appRes, err := l.svcCtx.PlatformRpc.GetAgentApp(ctx, &platform_rpc.GetAgentAppReq{AppId: cfgRes.AppId})
		if err != nil || appRes == nil || appRes.App == nil {
			l.logger.Error(model.LogMsg{Text: "GetAgentApp 失败", Data: map[string]any{"digitalId": digitalID, "appId": cfgRes.AppId, "err": fmt.Sprintf("%v", err)}})
			return nil
		}
		app := appRes.App
		systemPrompt := app.SystemPrompt
		llmModelID := app.LlmModelId

		// 拼最近聊天记录上下文。
		historyCtx := l.buildHistoryContext(ctx, conversationID, userID, digitalID)

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

		// 调用 Python agent_rpc SendMessage（fire-and-forget：只取 code，不等终稿）。
		reply, err := l.svcCtx.AgentRpc.SendMessage(ctx, &agent.SendMessageReq{
			UserId:         userID,
			ConversationId: conversationID,
			ModelId:        llmModelID,
			PartsJson:      string(partsJSON),
			SystemPrompt:   systemPrompt,
			DigitalId:      digitalID,
			ClientMsgId:    messageID,
		})
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "AgentRpc.SendMessage 调用失败", Data: map[string]any{"digitalId": digitalID, "messageId": messageID, "err": err.Error()}})
			return nil
		}
		if reply.GetCode() != 0 {
			l.logger.Error(model.LogMsg{Text: "AgentRpc.SendMessage 返回非0", Data: map[string]any{"digitalId": digitalID, "messageId": messageID, "code": reply.GetCode(), "msg": reply.GetMsg()}})
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
func (l *AgentTaskConsumerLogic) buildHistoryContext(ctx context.Context, conversationID, userID, digitalID string) string {
	var msgs []chat_models.ChatMessage
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
		if m.SendUserID != nil && *m.SendUserID == digitalID {
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

// extractPayload 从 SendMessage 投递的信封 {payload: {...}} 中取出内层 payload map。
func extractPayload(envelope map[string]interface{}) (map[string]interface{}, error) {
	v, ok := envelope["payload"]
	if !ok || v == nil {
		return nil, fmt.Errorf("信封缺少 payload 字段")
	}
	inner, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("payload 不是 map[string]interface{}")
	}
	return inner, nil
}
