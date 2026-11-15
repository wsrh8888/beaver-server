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
	"beaver/commonconfig"
	"beaver/middleware/jwt"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd string

	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	FriendRpc       zrpc.RpcClientConf
	GroupRpc        zrpc.RpcClientConf
	UserRpc         zrpc.RpcClientConf
	ChatRpc         zrpc.RpcClientConf
	EmojiRpc        zrpc.RpcClientConf
	DatasyncRpc     zrpc.RpcClientConf
	NotificationRpc zrpc.RpcClientConf
	AuthRpc         zrpc.RpcClientConf
	FileRpc         zrpc.RpcClientConf
	PlatformRpc     zrpc.RpcClientConf
	// Agent beaver-agent(py, beaver-ai 业务线) S2S 出站端点；跨仓禁 gRPC，静态 X-Internal-Token。
	Agent commonconfig.AgentEndpoint
	// Circle 圈子服务（beaver-server-social/app/circle）S2S 端点：client_credentials 换票 + Bearer。
	Circle commonconfig.Endpoint
	// S2S 入站凭据：收费模块（beaver-social 等）经此换票后调本面 /internal/* 能力
	// （User/Chat/Notification/Friend 的内部编排端点）。机制同收费模块侧，fail-closed。
	S2S struct {
		Secret  string                     `json:",optional"`
		Clients map[string]jwt.ClientEntry `json:",optional"`
	}
}
