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
	"fmt"
	"strings"

	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/open/openevent"
	"beaver/common/const/mqwsconst"
	"beaver/common/models/ctype"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"beaver/utils/conversation"
)

type digitalEventPusher struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func newDigitalEventPusher(ctx context.Context, svcCtx *svc.ServiceContext) *digitalEventPusher {
	return &digitalEventPusher{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("digital_event_push", ctx)}
}

func (p *digitalEventPusher) tryPush(in *chat_rpc.SendMsgReq, msg *ctype.Msg) {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Error(model.LogMsg{Text: "数字伙伴事件推送panic", Data: map[string]any{"err": fmt.Sprintf("%v", r)}})
		}
	}()

	// 数字伙伴自己的回复不再投回去，避免群里来回叫醒。
	if isOpenDigitalSender(in.DeviceId) {
		return
	}

	conversationType, userIDs := conversation.ParseConversationWithType(in.ConversationId)
	eventBody := p.buildEventBody(in, msg)

	switch conversationType {
	case 1:
		p.pushPrivateChat(in.UserId, userIDs, eventBody)
	case 2:
		p.pushGroupAt(in.UserId, in.ConversationId, msg, eventBody)
	}
}

func (p *digitalEventPusher) pushPrivateChat(senderID string, userIDs []string, event map[string]interface{}) {
	if len(userIDs) != 2 {
		return
	}
	var peerID string
	for _, uid := range userIDs {
		if uid != senderID {
			peerID = uid
			break
		}
	}
	if peerID == "" {
		return
	}

	res, err := p.svcCtx.OpenRpc.GetDigitalByUserID(p.ctx, &open_rpc.GetDigitalByUserIDReq{DigitalUserId: peerID})
	if err != nil || res == nil || !res.Found {
		return
	}

	// builtin 模式：投 agent_task，由 beaver-agent 消费跑 LLM。
	if res.Runtime == "builtin" {
		p.pushAgentTask(res.DigitalUserId, event)
		return
	}

	// workflow 模式：投 workflow_task，由 workflow 引擎消费（按 Nodes/Edges 调度）。
	if res.Runtime == "workflow" {
		p.pushWorkflowTask(res.DigitalUserId, event)
		return
	}

	// external 模式：走 openevent → WS 长连接下发外部 agent。
	p.dispatch(res.DigitalUserId, openevent.EventIMMessageReceive, event)
}

func (p *digitalEventPusher) pushGroupAt(senderID, conversationID string, msg *ctype.Msg, event map[string]interface{}) {
	if msg == nil || len(msg.AtUserIDs) == 0 {
		return
	}

	for _, atUserID := range msg.AtUserIDs {
		res, err := p.svcCtx.OpenRpc.GetDigitalByUserID(p.ctx, &open_rpc.GetDigitalByUserIDReq{DigitalUserId: atUserID})
		if err != nil || res == nil || !res.Found {
			continue
		}
		groupEvent := copyEventMap(event)
		groupEvent["group_id"] = conversation.GetTargetIDByConversation(conversationID, senderID)
		// builtin 模式：投 agent_task。
		if res.Runtime == "builtin" {
			p.pushAgentTask(res.DigitalUserId, groupEvent)
			continue
		}
		// workflow 模式：投 workflow_task。
		if res.Runtime == "workflow" {
			p.pushWorkflowTask(res.DigitalUserId, groupEvent)
			continue
		}
		// external 模式：走 openevent。
		p.dispatch(res.DigitalUserId, openevent.EventIMMessageReceiveGroup, groupEvent)
	}
}

// dispatch external 模式：走 openevent → WS 长连接下发外部 agent。
func (p *digitalEventPusher) dispatch(digitalID, eventType string, event map[string]interface{}) {
	conversationID, _ := event["conversation_id"].(string)
	if err := openevent.Push(p.ctx, p.svcCtx.RocketMQ, digitalID, eventType, conversationID, event); err != nil {
		p.logger.Error(model.LogMsg{
			Text: "数字伙伴事件投递失败",
			Data: map[string]any{"digitalId": digitalID, "event": eventType, "err": err.Error()},
		})
	}
}

// pushAgentTask builtin 模式：投 agent_task topic，由 beaver-agent 消费跑 LLM。
// 消息体对齐 agent consumer.py 的 payload 结构：digitalId + userId + conversationId + content + messageId。
func (p *digitalEventPusher) pushAgentTask(digitalID string, event map[string]interface{}) {
	conversationID, _ := event["conversation_id"].(string)
	senderID, _ := event["sender_id"].(string)
	messageID, _ := event["message_id"].(string)
	content, _ := event["content"].(string)

	payload := map[string]interface{}{
		"messageId":      messageID,
		"userId":        senderID,
		"conversationId": conversationID,
		"content":       content,
		"digitalId":     digitalID,
	}
	envelope := map[string]interface{}{"payload": payload}
	if err := p.svcCtx.RocketMQ.SendMessage(p.ctx, mqwsconst.MqTopicAgentTask, envelope); err != nil {
		p.logger.Error(model.LogMsg{
			Text: "agent_task 投递失败",
			Data: map[string]any{"digitalId": digitalID, "messageId": messageID, "err": err.Error()},
		})
	}
}

// pushWorkflowTask workflow 模式：投 workflow_task topic，由 workflow 引擎消费（按 Nodes/Edges 调度）。
func (p *digitalEventPusher) pushWorkflowTask(digitalID string, event map[string]interface{}) {
	conversationID, _ := event["conversation_id"].(string)
	senderID, _ := event["sender_id"].(string)
	messageID, _ := event["message_id"].(string)
	content, _ := event["content"].(string)

	payload := map[string]interface{}{
		"messageId":      messageID,
		"userId":         senderID,
		"conversationId": conversationID,
		"content":       content,
		"digitalId":     digitalID,
	}
	envelope := map[string]interface{}{"payload": payload}
	if err := p.svcCtx.RocketMQ.SendMessage(p.ctx, mqwsconst.MqTopicWorkflowTask, envelope); err != nil {
		p.logger.Error(model.LogMsg{
			Text: "workflow_task 投递失败",
			Data: map[string]any{"digitalId": digitalID, "messageId": messageID, "err": err.Error()},
		})
	}
}

func (p *digitalEventPusher) buildEventBody(in *chat_rpc.SendMsgReq, msg *ctype.Msg) map[string]interface{} {
	event := map[string]interface{}{
		"sender_id":       in.UserId,
		"conversation_id": in.ConversationId,
		"message_id":      in.MessageId,
		"msg_type":        msgTypeName(msg),
		"content":         extractTextContent(msg),
		"mentions":        msg.AtUserIDs,
	}
	if event["mentions"] == nil {
		event["mentions"] = []string{}
	}
	return event
}

func msgTypeName(msg *ctype.Msg) string {
	if msg == nil {
		return "text"
	}
	switch msg.Type {
	case ctype.TextMsgType:
		return "text"
	case ctype.MarkdownMsgType:
		return "markdown"
	case ctype.ImageMsgType:
		return "image"
	default:
		return "unknown"
	}
}

func extractTextContent(msg *ctype.Msg) string {
	if msg == nil {
		return ""
	}
	if msg.ReplyMsg != nil && msg.ReplyMsg.ReplyMsg != nil {
		if text := extractTextContent(msg.ReplyMsg.ReplyMsg); text != "" {
			return text
		}
	}
	if msg.TextMsg != nil {
		return msg.TextMsg.Content
	}
	if msg.MarkdownMsg != nil {
		return msg.MarkdownMsg.Content
	}
	return ""
}

func copyEventMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func isDigitalPeer(ctx context.Context, openRpc open_rpc.OpenClient, peerUserID string) bool {
	if peerUserID == "" {
		return false
	}
	res, err := openRpc.GetDigitalByUserID(ctx, &open_rpc.GetDigitalByUserIDReq{DigitalUserId: peerUserID})
	return err == nil && res != nil && res.Found
}

func isOpenDigitalSender(deviceID string) bool {
	return strings.EqualFold(deviceID, "open_digital")
}
