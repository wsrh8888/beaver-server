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

// Package agent beaver-agent（py，beaver-ai 业务线）的 S2S 出站客户端。
//
// beaver-agent 与本仓是跨团队边界：禁止 etcd gRPC 直连（团队铁律），统一走对方
// 标准机器视图 /beaver-ai/api/agent/chat/v1/* 端点。鉴权与 beaver-ai 的其它机器面一致：
// client_credentials 换票 + Bearer（s2s.clients 按 client_id 白名单，本仓登记为 beaver-server）。
// 域名由 commonconfig.Load 从 yaml 的 RemoteProxy.Agent.BaseURL 灌入 commonconfig.AgentUrl。
// 业务侧直接调包函数 SendMessageApi(...) / ListConversationsApi(...)。
package agent

import (
	"sync"
	"time"
)

// client 本包唯一的出站客户端，Init 时装配。
var client *Client

// Init 装配 beaver-ai 分配给本仓的 client_credentials（与 RemoteProxy.Ai 同凭据，
// 都是 beaver-ai 同一个服务），换票协议在 token.go。只在该服务装配依赖时调用一次。
func Init(clientID, clientSecret string) {
	client = &Client{clientID: clientID, clientSecret: clientSecret}
}

// Client 持凭据与换来的票，换票缓存挂在它上面。
type Client struct {
	clientID     string
	clientSecret string

	mu        sync.Mutex
	bearer    string
	expiresAt time.Time
}
