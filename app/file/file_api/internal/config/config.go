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
	Etcd  string
	Redis struct {
		Addr     string
		Password string
		Db       int
	}
	Domain      string // 项目对外访问域名（用于生成本地文件完整URL）
	FileMaxSize map[string]float64
	WhiteList   []string
	BlackList   []string
	UserRpc     zrpc.RpcClientConf
	Local       struct {
		UploadDir   string // 本地文件上传目录
		ProjectName string // 项目名称，用于文件路径前缀（为空则使用根目录）
	}
	Qiniu struct {
		ProjectName string // 项目名称，用于文件路径前缀（为空则使用根目录）
		AK          string
		SK          string
		Bucket      string
		Domain      string // 七牛云文件访问域名
		ExpireTime  int64
	}
}
