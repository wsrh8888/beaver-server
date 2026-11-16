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

package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"beaver/commonconfig"
	agenttype "beaver/types/remote/agent"
)

// SendMessage 发一句话给 Agent（fire-and-forget：只取 code，不等终稿）。
func SendMessage(ctx context.Context, header http.Header, body *agenttype.SendMessageReq) (*agenttype.Reply, error) {
	return post(ctx, "/beaver-ai/api/agent/chat/v1/send_message", body)
}

// ListConversations 桌面 datasync：增量列会话（since = 本地游标 version）。
func ListConversations(ctx context.Context, header http.Header, body *agenttype.ListConversationsReq) (*agenttype.Reply, error) {
	return post(ctx, "/beaver-ai/api/agent/chat/v1/list_conversations", body)
}

// post 发一次 POST 到 agent 标准机器视图；静态 X-Internal-Token 鉴权，响应透传 Reply{code,msg,resultJson}。
// agent 不用 go-zero 信封，故不走 ajax.Post（它按 {code,msg,result} 解包），这里自管传输与解码。
func post(ctx context.Context, path string, req any) (*agenttype.Reply, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("agentapi %s: marshal req: %w", path, err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, commonconfig.AgentUrl+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("agentapi %s: build request: %w", path, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if client != nil && client.token != "" {
		httpReq.Header.Set("X-Internal-Token", client.token)
	}
	resp, err := client.hc.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("agentapi %s: %w", path, err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("agentapi %s: read body: %w", path, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("agentapi %s: http %d: %s", path, resp.StatusCode, truncate(data))
	}
	var reply agenttype.Reply
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, fmt.Errorf("agentapi %s: decode reply: %w", path, err)
	}
	return &reply, nil
}

func truncate(b []byte) string {
	const max = 512
	if len(b) > max {
		return string(b[:max]) + "...(truncated)"
	}
	return string(b)
}
