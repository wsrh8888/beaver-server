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
// 网关 agent_api（FastAPI :22300）的 /internal/agent.<Method> 端点，静态
// X-Internal-Token 鉴权（与 agent_api etc/agent.yaml 的 internal_token 同值）。
// 结构对齐 remote/{open,ai,circle}：Init 灌入令牌并把域名写入 commonconfig.AgentUrl；
// 业务侧直接调包函数 agent.SendMessage(...) / agent.ListConversations(...)。
package agent

import (
	"net/http"
	"time"

	"beaver/commonconfig"
	"beaver/utils/ajax"
)

// agentTimeout 单次调用超时。agent 是同步短请求，10s 足够，不做成可配项。
const agentTimeout = 10 * time.Second

// client 本包唯一的出站客户端，Init 时装配。
var client *Client

// Init 装配 beaver-agent 的静态 X-Internal-Token 与网关地址，并把域名写入 commonconfig.AgentUrl。
// 只在该服务装配依赖时调用一次。
func Init(baseURL, internalToken string) {
	commonconfig.AgentUrl = baseURL
	client = &Client{
		token: internalToken,
		hc:    ajax.HTTPClient(agentTimeout),
	}
}

// Client 持静态令牌与共享连接池的 http 客户端。
type Client struct {
	token string
	hc    *http.Client
}
