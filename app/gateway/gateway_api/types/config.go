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

// Config 网关负责：验外部票 → 写可信身份头 → 注入机器票 → 转发。
// 业务服务默认只验 S2S，身份由网关写入的头传递，故网关持有验票密钥与机器身份。
//
// 三类外部票的密钥各自独立成段（AuthBeaver / AuthOpen / AuthAdmin）：
// 现阶段三份可填同一密钥，但配置形状必须分开，将来轮换其中一类不影响其余两类。
// 机器票（Machine）不混进外部票段——它不是「谁在用这个系统」，而是「网关以什么身份调下游」。
type Config struct {
	Name string
	Addr string `json:",optional"`
	Etcd string
	// ServiceAlias URL 前缀 → 实际服务名（服务合并后按旧前缀路由到合并服务，
	// 如 friend→relation；未命中映射时按 {前缀}_api 默认规则）。
	ServiceAlias map[string]string `json:",optional"`
	Log          logx.LogConf
	Limit        LimitConfig

	// 三类外部票验签密钥，按平面各占一段，互不影响。
	AuthBeaver AuthConfig `json:",optional"`
	AuthOpen   AuthConfig `json:",optional"`
	AuthAdmin  AuthConfig `json:",optional"`
	// Machine 网关作为内部调用方注入到业务的机器票（业务默认只挂 JwtS2s）。
	Machine MachineConfig `json:",optional"`
	// PublicPaths 不验外部票的路径前缀白名单（登录、换票口、文件预览等）。
	PublicPaths []string `json:",optional"`
	// AuthGroups 路径前缀 → 外部票类型（beaver/open/admin/public），未命中默认 beaver。
	AuthGroups map[string]string `json:",optional"`
}

// AuthConfig 单个平面的外部票验签密钥。
type AuthConfig struct {
	AccessSecret string `json:",optional"`
}

// MachineConfig 网关注入到下游业务的机器票（S2S）配置。
// 密钥按目标团队拆分：一把密钥对应一个受众/团队，签发方与验签方必须用同一把。
// 网关按请求路径的目标团队选 Secrets[team] 签发，下游服务用该团队自己的 S2S.Secret 验签。
// 查表顺序：Secrets[team/surface] → Secrets[team] → DefaultSecret → Secret（最后一项兼容旧配置）。
// Audiences 同理按 team/surface → team → DefaultAudience 选票据受众。
type MachineConfig struct {
	Secrets         map[string]string `json:",optional"`
	DefaultSecret   string            `json:",optional"`
	Secret          string            `json:",optional"` // 旧配置兜底，新配置用 Secrets/DefaultSecret
	DefaultAudience string            `json:",optional"`
	Audiences       map[string]string `json:",optional"`
	ClientID        string            `json:",optional"`
	Issuer          string            `json:",optional"`
}

type LimitConfig struct {
	Enable bool    `json:",default=false"`
	Rate   float64 `json:",default=100"`
	Burst  int     `json:",default=200"`
}
