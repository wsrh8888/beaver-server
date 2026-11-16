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

// Package ai AI 服务（beaver-ai 独立仓）的 S2S 出站客户端。
//
// AI 服务已从本仓拆出（platform_rpc 的 AI 面随 beaver-ai 独立部署）；跨仓不允许 etcd gRPC
// 直连，本包经 beaver-ai 标准机器视图 /beaver-ai/api/ai/*/v1/* 换票调用。Req/Res 放到 types/remote/ai，
// 本仓不依赖 beaver-ai 模块。
//
// 结构对齐 remote/circle（enter/remote/token + 业务方法文件）：Init 灌入凭据并把域名写入
// commonconfig.AiUrl；业务侧直接调包函数 GetAgent(...)，不在调用行里重复写地址与换票。
package ai

import (
	"sync"
	"time"
)

// client 本包唯一的出站客户端，Init 时装配。
var client *Client

// Init 装配出站凭据（beaver-ai 分配给本仓的 client_credentials），并把域名写入 commonconfig。
// 只在该服务装配依赖时调用一次；域名来自 Init 的 baseURL（与 circle/doc/moment 一致，
// 后续 Phase 4 统一收口到 commonconfig.Load，Init 去 baseURL）。
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
