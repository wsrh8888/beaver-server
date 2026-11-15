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
package agent

// Reply 对齐 agent_rpc.proto 的 Reply（proto JSON：result_json → resultJson）。
type Reply struct {
	Code       int32  `json:"code"`
	Msg        string `json:"msg"`
	ResultJson string `json:"resultJson"`
}

func (r *Reply) GetCode() int32 {
	if r == nil {
		return 0
	}
	return r.Code
}

func (r *Reply) GetMsg() string {
	if r == nil {
		return ""
	}
	return r.Msg
}

// SendMessageReq 对齐 agent_rpc.proto 的 SendMessageReq。
type SendMessageReq struct {
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
	Module         string `json:"module"`
	AppId          string `json:"appId"`
}

// ListConversationsReq 对齐 agent_rpc.proto 的 ListConversationsReq（since: version 游标）。
type ListConversationsReq struct {
	UserId string `json:"userId"`
	Since  int64  `json:"since"`
	Module string `json:"module"`
	Limit  int32  `json:"limit"`
}
