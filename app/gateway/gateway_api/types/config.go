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
	// ServiceAlias URL 前缀 → 实际服务名（服务合并后按旧前缀路由到合并服务，
	// 如 friend→relation；未命中映射时按 {前缀}_api 默认规则）。
	ServiceAlias    map[string]string `json:",optional"`
	Log        logx.LogConf
	Prometheus PrometheusConfig
	Limit          LimitConfig
	Auth           AuthConfig
	// S2S 队间机器票验签配置（与 account_api 签发密钥同源）。
	S2S            S2SConfig
	CompatibleVersion string   `json:",optional"` // 已停用，网关不再按客户端版本拦截
	PublicList        []string `json:",optional"` // Gateway 不鉴权（含 *_public）
	CustomAuthList    []string `json:",optional"` // 透传，由下游服务 middleware 鉴权
}

type AuthConfig struct {
	AccessSecret string
	AccessExpire int
}

// S2SConfig 网关验 S2S 机器票的密钥与受众。
//
// Secret 必须与 account_api 签发 /api/auth/token 的 S2S.Secret 一致；
// Audience 留空时取 authgate.AudienceCore（core:api）。
type S2SConfig struct {
	Secret   string
	Audience string `json:",optional"`
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
