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

// Package authgate 是「用户票或机器票」共用的身份闸门。
//
// 声明式终态：本闸门只由 .api 的 `@server middleware: AuthGate` 挂载，不再由
// main 里的 server.Use 全局接管。公开组不挂任何鉴权声明，用户组用 go-zero 内置
// `jwt: Auth`，需要「用户或机器」的组用本闸门。
//
// 本闸门只做 Authentication（验签 + 写 context 身份结论）；Authorization（是否
// 允许机器调用、scopes 是否足够）一律留给 logic。
//
// fail-closed：两条通道都没验过一律 401，不存在「没配密钥就放行」的分支；
// 外部传入的身份头一律先删再按验签结论重建，杜绝伪造。
package authgate

import (
	"net/http"
	"strings"

	"beaver/common/identity"
	"beaver/middleware/jwt"
	"beaver/utils/jwts"
)

const (
	// HeaderUserID 已验签的用户身份头。仅作下游透传展示，logic 不得据此判定身份。
	HeaderUserID = "Beaver-User-Id"
	// HeaderClientID 已验签的调用方标识头。
	HeaderClientID = jwt.HeaderClientID
	// HeaderToken 历史用户令牌头名；本服务直接校验其 JWT，等价于 Authorization: Bearer。
	HeaderToken = "Token"
	// AudienceCore 本团队（A 主体）全部业务面共用的机器票受众。
	AudienceCore = "core:api"
)

const bearerPrefix = "Bearer "

// Gate 闸门本体：持两个密钥与受众，负责把一次请求判成用户 / 机器 / 拒绝。
type Gate struct {
	userSecret      string
	s2sSecret       string
	s2sAud          string
	allowQueryToken bool
}

// Option 闸门构造选项。
type Option func(*Gate)

// WithQueryToken 允许从 URL query 的 token 参数取用户票（WebSocket 握手无法带自定义头）。
func WithQueryToken() Option {
	return func(g *Gate) { g.allowQueryToken = true }
}

// NewGate 构造闸门。密钥缺失由 Validate 在启动期拦截。
func NewGate(userSecret, s2sSecret, s2sAud string, opts ...Option) *Gate {
	g := &Gate{
		userSecret: userSecret,
		s2sSecret:  s2sSecret,
		s2sAud:     s2sAud,
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Validate 启动期自检：用户密钥、S2S 密钥与受众都必须配齐，缺一即拒绝启动。
func (g *Gate) Validate() error {
	if strings.TrimSpace(g.userSecret) == "" {
		return errMissingUserSecret
	}
	if strings.TrimSpace(g.s2sSecret) == "" {
		return errMissingS2SSecret
	}
	if strings.TrimSpace(g.s2sAud) == "" {
		return errMissingS2SAud
	}
	return nil
}

// Resolve 判定一次请求的身份。
//
// 顺序：S2S 机器票 → 用户 JWT → 拒绝。先判机器票是因为两者都走 Authorization
// 头，机器票的 aud 与用户票不同，验不过会自然落到用户通道。
// 公开组不挂本闸门，因此这里不再有任何路径白名单。
func (g *Gate) Resolve(r *http.Request) (identity.Identity, bool) {
	bearer := bearerOf(r)

	// 1) 队间机器票。
	if bearer != "" {
		if clientID, scopes, err := jwt.Verify(g.s2sSecret, g.s2sAud, bearer); err == nil {
			return identity.Identity{
				Kind:     identity.KindS2S,
				ClientID: clientID,
				Scopes:   scopes,
			}, true
		}
	}

	// 2) 用户 JWT：统一走 Authorization: Bearer；WebSocket 握手可退回 query token。
	userToken := bearer
	if userToken == "" && g.allowQueryToken {
		userToken = strings.TrimSpace(r.URL.Query().Get("token"))
	}
	if userToken != "" {
		if claims, err := jwts.ParseToken(userToken, g.userSecret); err == nil && claims.UserID != "" {
			return identity.Identity{Kind: identity.KindUser, UserID: claims.UserID}, true
		}
	}

	return identity.Identity{}, false
}

// Apply 把身份结论写回请求，并把结论注入 context。
//
// 无论哪种身份，都先删掉外部可能伪造的两个身份头，再按验签结论重建；
// logic 只读 context（common/identity.From），不再依赖任何请求头。
func (g *Gate) Apply(r *http.Request, id identity.Identity) *http.Request {
	r.Header.Del(HeaderClientID)
	r.Header.Del(HeaderUserID)

	switch id.Kind {
	case identity.KindUser:
		r.Header.Set(HeaderUserID, id.UserID)
	case identity.KindS2S:
		r.Header.Set(HeaderClientID, id.ClientID)
	}

	return r.WithContext(identity.With(r.Context(), id))
}

// bearerOf 取 Authorization 头里的裸令牌（无则空串）。
func bearerOf(r *http.Request) string {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(raw, bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(raw, bearerPrefix))
}
