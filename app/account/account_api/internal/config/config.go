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
	// Open 开放接口服务（beaver-open 独立仓）S2S 端点：client_credentials 换票 + Bearer。
	Open commonconfig.Endpoint
	Auth        struct {
		AccessSecret string
		AccessExpire int
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
