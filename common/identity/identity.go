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

// Package identity 是全仓唯一的身份结论载体。
//
// Authentication 由 .api @server 声明的中间件完成，中间件把验签结论写进 context；
// Authorization 留在 logic，logic 只能通过本包从 context 读取身份，不再信任请求头或
// 请求体里由调用方自带的 userId/clientId。
package identity

import "context"

// Kind 身份来源。
type Kind int

const (
	// KindAnonymous 未验签（公开组，或凭据缺失）。
	KindAnonymous Kind = iota
	// KindUser 终端用户（用户 JWT 验签通过）。
	KindUser
	// KindS2S 队间机器（client_credentials 机器票验签通过）。
	KindS2S
)

// Identity 一次请求的身份结论，由鉴权中间件写入 context。
type Identity struct {
	Kind     Kind
	UserID   string   // Kind==KindUser 时为验签得出的用户 ID
	ClientID string   // Kind==KindS2S 时为票据里的调用方 client_id
	Scopes   []string // Kind==KindS2S 时为票据携带的授权域
}

func (i Identity) IsUser() bool    { return i.Kind == KindUser }
func (i Identity) IsMachine() bool { return i.Kind == KindS2S }

type identityCtxKey struct{}

// With 把身份结论注入 context。只应由鉴权中间件调用。
func With(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, identityCtxKey{}, id)
}

// From 读取本次请求的身份结论：只认鉴权中间件通过 With 注入的结构化值。
func From(ctx context.Context) Identity {
	if ctx == nil {
		return Identity{}
	}
	if v, ok := ctx.Value(identityCtxKey{}).(Identity); ok {
		return v
	}
	return Identity{}
}

// UserID 取验签得出的用户 ID；无用户身份时为空串。
func UserID(ctx context.Context) string { return From(ctx).UserID }

// ClientID 取验签得出的调用方 client_id；非机器调用时为空串。
func ClientID(ctx context.Context) string { return From(ctx).ClientID }

// IsMachine 是否队间机器调用，供 logic 做 Authorization 判断。
func IsMachine(ctx context.Context) bool { return From(ctx).IsMachine() }

// RequireMachine 纯队间能力面的 Authorization 结论：不是机器票即拒绝。
func RequireMachine(ctx context.Context) bool { return From(ctx).IsMachine() }
