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
	"beaver/middleware/jwt"

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
	// RemoteProxy 出站 S2S 端点注册表：key 是对端团队/服务名，值是对方网关地址与鉴权方式。
	// 跨团队一律走对方网关的 /internal/<Module>.<Method>，禁止 etcd gRPC 直连、禁止引用对方代码。
	RemoteProxy commonconfig.RemoteProxy `json:",optional"`
	Auth        struct {
		AccessSecret string
		AccessExpire int
	}
	// S2S 入站机器票：本团队业务面共用一套凭据，票据受众统一 core:api。
	// 换票端点在本服务的 POST /api/auth/token（client_credentials）。
	S2S struct {
		Secret  string                     `json:",optional"`
		Clients map[string]jwt.ClientEntry `json:",optional"`
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
