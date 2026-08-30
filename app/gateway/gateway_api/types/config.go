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

package types

import "github.com/zeromicro/go-zero/core/logx"

type Config struct {
	Name       string
	Addr       string `json:",optional"`
	Etcd       string
	Log        logx.LogConf
	Prometheus PrometheusConfig
	Limit          LimitConfig
	Auth           AuthConfig
	PublicList     []string `json:",optional"` // Gateway 不鉴权（含 *_public）
	CustomAuthList []string `json:",optional"` // 透传，由下游服务 middleware 鉴权
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int
}

type PrometheusConfig struct {
	Enable bool   `json:",default=false"`
	Path   string `json:",default=/metrics"`
}

type LimitConfig struct {
	Enable bool    `json:",default=false"`
	Rate   float64 `json:",default=100"`
	Burst  int     `json:",default=200"`
}
