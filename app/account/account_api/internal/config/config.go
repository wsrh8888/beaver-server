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
 *
 * beaver-server-header-v1
 */

package config

import (
	"beaver/commonconfig"
	"beaver/utils/s2s"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd  string
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	RocketMQ struct {
		Addr string
	}
	WhiteList []string
	// AccountRpc 本域 rpc（user + auth 同进程，etcd accountrpc.rpc）
	AccountRpc      zrpc.RpcClientConf
	RelationRpc     zrpc.RpcClientConf // relation_rpc（friend+group）
	NotificationRpc zrpc.RpcClientConf
	// RemoteProxy 唯一出站注册表：本仓出站到其它团队走各自标准 HTTP 契约 + 机器票 Bearer；
	// 外仓打 A 主体（Beaver）统一走 gateway_api 的 /beaver/api/... 标准契约，业务能力只有这一套 HTTP 契约；
	// 本仓出站不使用任何 /internal/* 私有端点，只走对端标准 .api 路径。
	RemoteProxy commonconfig.RemoteProxy `json:",optional"`
	Auth        struct {
		AccessSecret string
		AccessExpire int
	}
	// S2S 入站机器票：本团队业务面共用一套凭据，票据受众统一 core:api。
	// 换票端点在本服务的 POST /beaver/api/account/auth_token/v1/issue（client_credentials）。
	S2S struct {
		Secret  string                     `json:",optional"`
		Clients map[string]s2s.ClientEntry `json:",optional"`
	}
	Email struct {
		QQ struct {
			Host     string
			Port     int
			Username string
			Password string
		}
	}
}
