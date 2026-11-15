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

// Package authgate 是业务 HTTP 面的唯一身份闸门：同一套路由既服务终端用户，也服务队间机器调用。
//
// 一条业务能力只有一套 .api 契约（例如 /api/user/v1/list_info）。差异只在「谁在调用」：
//   - 终端用户：网关验过用户 JWT 后仍把 Token 头转发下来，本闸门再验一次签名，
//     得出 UserID 并覆盖 Beaver-User-Id —— 不信任任何外部直接塞进来的身份头；
//   - 队间机器：Authorization: Bearer <client_credentials 换来的 S2S 票>，验签得出 ClientID，
//     并可携带 Beaver-User-Id 声明「代表哪个用户行事」（机器代办场景，如管理端批量操作）。
//
// logic 只读 Beaver-User-Id（既有 header tag 约定不变），需要区分来源时读 context 里的 Identity。
//
// fail-closed：两条通道都没验过一律 401，不存在「没配密钥就放行」的分支。
package authgate

import (
	"context"
	"net/http"
	"strings"

	"beaver/middleware/jwt"
	"beaver/utils/jwts"
)

const (
	// HeaderUserID 已验签的身份结论头：用户通道写验签得出的 UserID，
	// 机器通道透传调用方声明的被代表用户（可为空）。
	HeaderUserID = "Beaver-User-Id"
	// HeaderClientID 已验签的调用方标识头：机器通道写票据里的 client_id，
	// 用户通道留空（用户不是 client）。外部传入的同名头一律先剥离。
	HeaderClientID = jwt.HeaderClientID
	// HeaderToken 用户访问令牌头，与管理端 / 客户端约定一致。
	HeaderToken = "Token"
)

// AudienceCore 本团队（A 主体）全部业务面共用的机器票受众。
//
// 一个受众而不是每服务一个：收费仓换一次票就能打 account / message / relation /
// asset / notification / platform 六个面，省掉每服务各换一次票的往返。
const AudienceCore = "core:api"

// DefaultWhiteList 业务面默认免鉴权前缀，与网关 PublicList 对齐。
//
// 各服务按自己的路由表取用；未命中的前缀不会造成放行（只是不匹配）。
var DefaultWhiteList = []string{
	"/api/auth/auth_public/",       // 登录 / 注册 / 验证码 / 扫码
	"/api/auth/token",              // client_credentials 换 S2S 票
	"/api/platform/track_public/",  // 埋点上报
	"/api/platform/update_public/", // 版本检查
	"/api/platform/v1/instance",    // 实例公开信息
	"/api/file/preview/",           // 文件预览（fileKey 即凭据）
}

const bearerPrefix = "Bearer "

// Kind 身份来源。
type Kind int

const (
	// KindAnonymous 白名单路径，未做身份判定。
	KindAnonymous Kind = iota
	// KindUser 终端用户（用户 JWT 验签通过）。
	KindUser
	// KindS2S 队间机器（client_credentials 机器票验签通过）。
	KindS2S
)

// Identity 一次请求的身份结论，由闸门写入 context，logic 按需读取。
type Identity struct {
	Kind     Kind
	UserID   string   // Kind==KindUser 时为验签得出的用户 ID
	ClientID string   // Kind==KindS2S 时为票据里的调用方 client_id
	Scopes   []string // Kind==KindS2S 时为票据携带的授权域
}

// IsMachine 是否队间机器调用。
func (i Identity) IsMachine() bool { return i.Kind == KindS2S }

// IsUser 是否终端用户本人调用。
func (i Identity) IsUser() bool { return i.Kind == KindUser }

type identityCtxKey struct{}

// WithIdentity 把身份结论注入 context。
func WithIdentity(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, identityCtxKey{}, id)
}

// IdentityFromContext 取本次请求的身份结论（白名单路径为零值）。
func IdentityFromContext(ctx context.Context) Identity {
	if v, ok := ctx.Value(identityCtxKey{}).(Identity); ok {
		return v
	}
	return Identity{}
}

// Gate 闸门本体：持两个密钥与受众，负责把一次请求判成用户 / 机器 / 拒绝。
type Gate struct {
	userSecret string   // 用户 JWT 签名密钥（与网关 Auth.AccessSecret 同值）
	s2sSecret  string   // S2S 机器票签名密钥（本团队 S2S.Secret）
	s2sAud     string   // 机器票受众，本团队业务面统一为 core:api
	whiteList  []string // 免鉴权路径前缀
}

// NewGate 构造闸门。两个密钥都为空时构造失败由调用方在启动期拦截（见 Validate）。
func NewGate(userSecret, s2sSecret, s2sAud string, whiteList []string) *Gate {
	return &Gate{
		userSecret: userSecret,
		s2sSecret:  s2sSecret,
		s2sAud:     s2sAud,
		whiteList:  whiteList,
	}
}

// Validate 启动期自检：用户密钥与 S2S 密钥都必须配齐，缺一即拒绝启动。
//
// 不提供「只配一个也能跑」的降级——那样会静默放行另一半调用方。
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

// IsWhiteListed 路径是否免鉴权。
func (g *Gate) IsWhiteListed(path string) bool {
	for _, p := range g.whiteList {
		if p == "" {
			continue
		}
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

// Resolve 判定一次请求的身份。
//
// 顺序：白名单 → S2S 机器票 → 用户 JWT → 拒绝。先判机器票是因为两者都走 Authorization
// 头，机器票的 aud 与用户票不同，验不过会自然落到用户通道。
func (g *Gate) Resolve(r *http.Request) (Identity, bool) {
	if g.IsWhiteListed(r.URL.Path) {
		return Identity{Kind: KindAnonymous}, true
	}

	bearer := bearerOf(r)

	// 1) 队间机器票。
	if bearer != "" {
		if clientID, scopes, err := jwt.Verify(g.s2sSecret, g.s2sAud, bearer); err == nil {
			return Identity{Kind: KindS2S, ClientID: clientID, Scopes: scopes}, true
		}
	}

	// 2) 用户 JWT：优先 Token 头（网关与客户端约定），退化到 Authorization。
	userToken := strings.TrimSpace(r.Header.Get(HeaderToken))
	if userToken == "" {
		userToken = bearer
	}
	if userToken != "" {
		if claims, err := jwts.ParseToken(userToken, g.userSecret); err == nil && claims.UserID != "" {
			return Identity{Kind: KindUser, UserID: claims.UserID}, true
		}
	}

	return Identity{}, false
}

// Apply 把身份结论写回请求头，返回带 context 的请求。
//
// 头部处理是这一步的关键：
//   - 机器通道：剥离外部 X-Client-Id 后按验签结论重建；Beaver-User-Id 透传（机器代表用户），
//     但只有非空才写，避免把空值盖到下游。
//   - 用户通道：Beaver-User-Id 按验签结论覆盖，外部塞进来的伪造值直接失效。
func (g *Gate) Apply(r *http.Request, id Identity) *http.Request {
	r.Header.Del(HeaderClientID)

	switch id.Kind {
	case KindUser:
		r.Header.Set(HeaderUserID, id.UserID)
	case KindS2S:
		r.Header.Set(HeaderClientID, id.ClientID)
		if id.UserID != "" {
			r.Header.Set(HeaderUserID, id.UserID)
		}
	}

	return r.WithContext(WithIdentity(r.Context(), id))
}

// bearerOf 取 Authorization 头里的裸令牌（无则空串）。
func bearerOf(r *http.Request) string {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(raw, bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(raw, bearerPrefix))
}
