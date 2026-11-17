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

// Package agent beaver-agent S2S 出入参类型，与对端 agent_rpc proto JSON 逐字段对齐。
//
// 放这里而非 remote/agent 内部，便于调用方只依赖类型、不依赖客户端实现。
// 每接口一组 type(...) 放 Req + Res + 本接口专属子类型。
// 对端已按 go-zero 标准信封 {code,msg,result} 返回：code/msg 由 ajax 统一判成败，
// 各 Res 只描述 result 里的业务字段，与 open/ai 等 peer 同口径。
package agent

// SendMessageApi 发消息：Req + Res（result 无业务字段，空 Res）。
type (
	// SendMessageReq 对齐 agent_rpc.proto 的 SendMessageReq。
	SendMessageReq struct {
		UserId         string `json:"userId"`
		DeviceId       string `json:"deviceId"`
		TraceId        string `json:"traceId"`
		ConversationId string `json:"conversationId"`
		TopicId        string `json:"topicId"`
		ModelId        string `json:"modelId"`
		ModelSource    string `json:"modelSource"`
		// 多模态 parts JSON 数组（权威用户内容）。
		PartsJson string `json:"partsJson"`
		// 本机空间 id；空=云端轻会话，有值=本机空间。
		WorkspaceId    string `json:"workspaceId"`
		PermissionMode string `json:"permissionMode"`
		ClientMsgId    string `json:"clientMsgId"`
		WorkspacePath  string `json:"workspacePath"`
		TargetDeviceId string `json:"targetDeviceId"`
		SystemPrompt   string `json:"systemPrompt"`
		DigitalId      string `json:"digitalId"`
		// 数字伙伴绑定的 Agent；非空时 AI 侧按 agentId 取人设/模型/工具，SystemPrompt 传空。
		AgentId string `json:"agentId"`
		Module  string `json:"module"`
	}

	// SendMessageRes 发消息回包的 result（fire-and-forget，无业务字段）。
	SendMessageRes struct{}
)

// ListConversationsApi 列会话：Req + Res + 本接口专属子类型。
type (
	// ListConversationsReq 对齐 agent_rpc.proto 的 ListConversationsReq（since: version 游标）。
	ListConversationsReq struct {
		UserId string `json:"userId"`
		Since  int64  `json:"since"`
		Module string `json:"module"`
		Limit  int32  `json:"limit"`
	}

	// ListConversationsRes 列会话回包的 result。
	ListConversationsRes struct {
		List          []ListConversationsConversationItem `json:"list"`
		ServerVersion int64                               `json:"serverVersion"`
	}

	// ListConversationsConversationItem 会话版本摘要项（仅本接口使用）。
	ListConversationsConversationItem struct {
		ConversationID string `json:"conversationId"`
		Title          string `json:"title"`
		Avatar         string `json:"avatar"`
		Module         string `json:"module"`
		WorkspaceId    string `json:"workspaceId"`
		UpdatedAt      int64  `json:"updatedAt"`
		LatestTopicId  string `json:"latestTopicId"`
		Preview        string `json:"preview"`
		Version        int64  `json:"version"`
	}
)
