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

// Package commonconfig 出站 S2S 端点注册表（仓级共享，五个仓内容完全一致）。
//
// 团队之间不互通：本仓不对外提供 rpc，也不直连别人的 etcd；要用其它团队的能力，
// 一律走各自标准 HTTP 契约 + 机器票 Bearer，不跨仓 gRPC 直连、不经业务门面中转。
//
// 出站分两类：外仓打 A 主体（Beaver）统一走 gateway_api 的 /beaver/api/... 标准契约；
// 本仓出站到其它团队按上述方式。本仓出站不使用任何 /internal/* 私有端点。
//
// 两条轨合成一条：yaml 的 RemoteProxy 段是部署事实（唯一出站注册表），各服务 config
// 直接内嵌本包的 RemoteProxy 类型（不再各写一份字段表），启动时 Load 一次把域名灌进
// 包级变量；remote/<对端>.Init 只收凭据，域名一律从这里读。业务侧不碰本包，只调 remote
// 的包函数。Beaver 是 A 主体整体的 peer（唯一入口 gateway_api）；本仓自身不出站调自己，
// 这一项只在五仓之间保持类型对称，yaml 留空即可，禁止按业务域拆成多个 BaseURL。
//
// 只存原值：不去空白、不削末尾斜杠，拼地址的规则由调用方自己保证（多一个斜杠就是 404）。
package commonconfig

const (
	// AuthClientCredentials OAuth2 客户端凭据换票，跨团队出站的标准方式。
	AuthClientCredentials = "client_credentials"
)

// Endpoint 出站被调方 S2S 端点项：域名 + 鉴权方式 + 凭据。
//
// Auth 单独留成字段而不是隐含约定，是为了让 yaml 把「怎么鉴权」写清楚；
// 目前只有 client_credentials 一种合法值，为空按 client_credentials 处理。
type Endpoint struct {
	BaseURL      string `json:",optional"`
	Auth         string `json:",optional"`
	ClientID     string `json:",optional"`
	ClientSecret string `json:",optional"`
}

// RemoteProxy 出站 S2S 端点注册表：key 是对端团队/服务名，值是对方标准 .api 服务地址与凭据。
//
// 一个进程只跑一个服务，用不到的对端在 yaml 里留空即可（全部 optional），
// 所以五个仓共用同一份类型，不会出现「A 仓有这个字段、B 仓没有」的漂移。
type RemoteProxy struct {
	// Beaver A 主体 beaver-server 整体：唯一入口 gateway_api。
	// 本仓自用时留空（optional）；禁止把 account/message/relation/... 拆成多个 BaseURL。
	Beaver Endpoint `json:",optional"`
	// Ai AI 服务（beaver-ai）。
	Ai Endpoint `json:",optional"`
	// Open 开放平台（beaver-open）。
	Open Endpoint `json:",optional"`
	// Social 社交服务（beaver-social）：圈子与动态同属一个产品整体，只有一个入口。
	Social Endpoint `json:",optional"`
	// Doc 在线文档服务（beaver-server-online）。
	Doc Endpoint `json:",optional"`
	// Agent 智能体 py 服务（beaver-ai，/agent/chat/v1/*）：与 Ai 同凭据，client_credentials 换票。
	Agent Endpoint `json:",optional"`
}

// 各对端服务域名，由 Load 从服务自己的 yaml 灌入；该服务用不到的对端保持空。
var (
	// BeaverUrl A 主体 gateway_api 唯一入口；本仓不出站调自己，保持五仓命名对称。
	BeaverUrl string

	AiUrl     string
	OpenUrl   string
	SocialUrl string
	DocUrl    string
	AgentUrl  string
)

// Load 由各服务 NewServiceContext 在装配依赖前调一次，把 yaml 里配的域名存进本包。
// 一个进程只跑一个服务，不存在多份 Load 交叉覆盖。
func Load(p RemoteProxy) {
	BeaverUrl = p.Beaver.BaseURL

	AiUrl = p.Ai.BaseURL
	OpenUrl = p.Open.BaseURL
	SocialUrl = p.Social.BaseURL
	DocUrl = p.Doc.BaseURL
	AgentUrl = p.Agent.BaseURL
}
