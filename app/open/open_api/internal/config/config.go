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

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	Etcd    string
	UserRpc zrpc.RpcClientConf
	AuthRpc    zrpc.RpcClientConf
	ChatRpc    zrpc.RpcClientConf
	GroupRpc zrpc.RpcClientConf
	OpenRpc  zrpc.RpcClientConf

	// RocketMQ：接收 IM 事件（由 chat_rpc 投递）并推送给在线机器人连接
	RocketMQ struct {
		Addr string
	}

	// WebSocket：开放平台长连接参数
	WebSocket struct {
		PongWait       int // 读超时（秒）
		WriteWait      int // 写超时（秒）
		PingPeriod     int // 心跳间隔（秒），应小于 PongWait
		MaxMessageSize int // 单帧最大字节
	}
}
