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
	// RemoteProxy 出站 S2S 端点注册表：只登记其它团队（social / agent）的入口。
	// 本面只做增量同步，不承接 A 主体业务能力，也不对外提供业务门面。
	RemoteProxy commonconfig.RemoteProxy `json:",optional"`
}
