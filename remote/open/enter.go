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

// Package open A 主体对开放接口服务（beaver-open）的 S2S 出站客户端。
//
// 跨仓不允许 etcd gRPC 直连（团队铁律），统一走对方标准路径 /beaver-open/api/open/machine/v1/*（client_credentials
// 换票 + Bearer）。Req/Res 放到 types/remote/open（由 open_rpc proto 生成后整体迁入），
// 本包不依赖 beaver-open 模块。
//
// 结构对齐 remote/circle：Init 灌入凭据并把域名写入 commonconfig.OpenUrl；业务侧直接调包函数
// GetDigitalByUserID(...) / ExchangeToken(...)。
package open

import (
	"sync"
	"time"
)

// client 本包唯一的出站客户端，Init 时装配。
var client *Client

// Init 装配出站凭据（beaver-open 分配给本仓的 client_credentials），并把域名写入 commonconfig.OpenUrl。
// 只在该服务装配依赖时调用一次。
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
