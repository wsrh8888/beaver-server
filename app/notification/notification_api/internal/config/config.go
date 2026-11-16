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
	"beaver/middleware/jwt"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Etcd  string
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	RocketMQ struct {
		Addr string
	}
	UserRpc zrpc.RpcClientConf
	AuthRpc zrpc.RpcClientConf
	// NotificationRpc：本业务线的 gRPC 服务（etcd notificationrpc.rpc），队间能力面薄转发到它。
	NotificationRpc zrpc.RpcClientConf
	// Auth 用户 JWT 签名密钥，本服务独立签发并校验用户票：
	// 本服务独立验用户票，不信任任何外部直接塞进来的 Beaver-User-Id。
	Auth struct {
		AccessSecret string
		AccessExpire int
	}
	// S2S 入站机器票：队间调用与本服务共用同一套路由，凭据受众统一 core:api。
	S2S struct {
		Secret  string                     `json:",optional"`
		Clients map[string]jwt.ClientEntry `json:",optional"`
	}
}
