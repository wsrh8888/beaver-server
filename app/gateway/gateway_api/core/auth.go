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

package core

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"beaver/utils/jwts"
	"beaver/utils/s2s"
)

// 全仓统一身份头名。验票前一律剥离客户端自带值，验票后由网关按结论重建。
// user / open 平面写 Beaver-User-Id，admin 平面写 Beaver-Admin-Id——头名按平面分，
// 下游 .api 的字段名按域分（UserID / AdminID），两者互不牵连。
const (
	HeaderUserID      = "Beaver-User-Id"
	HeaderAdminID     = "Beaver-Admin-Id"
	HeaderPartnerID   = "X-Partner-Id"
	HeaderServiceName = "X-Service-Name"
	HeaderLegacyUser  = "X-User-Id"

	bearerPrefix = "Bearer "
)

// 客户端可伪造、网关必须先剥后写的身份头。
var spoofableHeaders = []string{
	HeaderUserID, HeaderAdminID, HeaderLegacyUser, HeaderPartnerID, HeaderServiceName,
}

// ticketKind 外部票平面。三类票密钥各自独立成段（见 jwtBeaver/jwtOpen/jwtAdmin.go），
// 受众也各自独立（beaver:user / beaver:developer / beaver:admin），互不通用。
type ticketKind int

const (
	ticketPublic ticketKind = iota
	ticketBeaver
	ticketOpen
	ticketAdmin
)

// stripIdentityHeaders 剥离所有客户端自带身份头（强制处理顺序第 1 步）。
func stripIdentityHeaders(r *http.Request) {
	for _, h := range spoofableHeaders {
		r.Header.Del(h)
	}
}

// routeAuthKind 按路径判定要验的外部票类型。
// 优先级：PublicPaths 精确白名单 > AuthGroups 最长前缀 > 默认主站用户票（beaver）。
// public 必须优先，否则 `/beaver-open/` 这类 team 级前缀会吞掉其下的登录/预览等公开端点。
func (p Proxy) routeAuthKind(path string) ticketKind {
	// 1) 公开白名单优先（最长匹配，允许精确到单个 action）。
	bestPublic := ""
	for _, prefix := range p.Config.PublicPaths {
		if strings.HasPrefix(path, prefix) && len(prefix) > len(bestPublic) {
			bestPublic = prefix
		}
	}
	if bestPublic != "" {
		return ticketPublic
	}

	// 2) team/surface 级鉴权分组，取最长前缀。
	best, kind := "", ticketBeaver
	for prefix, name := range p.Config.AuthGroups {
		if !strings.HasPrefix(path, prefix) || len(prefix) <= len(best) {
			continue
		}
		best = prefix
		switch strings.ToLower(name) {
		case "open":
			kind = ticketOpen
		case "admin":
			kind = ticketAdmin
		case "public":
			kind = ticketPublic
		default:
			kind = ticketBeaver
		}
	}
	return kind
}

// authenticate 强制处理顺序第 1~3 步：剥头 → 验票 → 写可信身份头。
//
// 验票分两条口：
//  1. 机器票优先：团队间调用带的是 S2S 票，受众就是目标团队，验过即放行，不再要求外部票。
//     机器面（machine/v1 等）本该如此，否则内部调用全被外部票闸门拦死。
//  2. 外部票：按路由平面分派到 jwtBeaver / jwtOpen / jwtAdmin 各自的验签实现。
//
// admin 平面写 Beaver-Admin-Id，其余平面写 Beaver-User-Id，让下游字段名与语义对齐。
// public 路由不验票，但同样先剥头，杜绝客户端自带身份头透传到下游。
// 返回 isMachine：true 表示请求携带的是内部机器票，此时不再注入网关票（保留原调用方身份）。
func (p Proxy) authenticate(r *http.Request) (bool, error) {
	stripIdentityHeaders(r)
	kind := p.routeAuthKind(r.URL.Path)
	if kind == ticketPublic {
		return false, nil
	}
	// 1) 机器票（服务间调用）：票面本身就是内部凭据，不再叠加外部票要求。
	if p.verifyMachineTicket(r) {
		return true, nil
	}
	// 2) 外部票：各平面各验各的密钥与受众。
	var (
		claims *jwts.CustomClaims
		err    error
	)
	switch kind {
	case ticketOpen:
		claims, err = p.verifyOpenTicket(r)
	case ticketAdmin:
		claims, err = p.verifyAdminTicket(r)
	default:
		claims, err = p.verifyBeaverTicket(r)
	}
	if err != nil {
		return false, err
	}
	r.Header.Set(identityHeader(kind), claims.UserID)
	return false, nil
}

// verifyTicket 三个平面共用的验签口径：取票 → 验签 → 校验非空与受众。
// 密钥与受众由调用方按平面传入，本函数不持有任何配置。
func verifyTicket(r *http.Request, secret, audience, missingSecretMsg string) (*jwts.CustomClaims, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, errors.New(missingSecretMsg)
	}
	token := ticketFromRequest(r)
	if token == "" {
		return nil, errors.New("缺少认证令牌")
	}
	claims, err := jwts.ParseToken(token, secret)
	if err != nil || claims.UserID == "" {
		return nil, errors.New("票据无效或已过期")
	}
	// 三票可能共用同一密钥，只有 Audience 分平面。缺了这一步，普通用户票能直接打后台接口。
	if !claims.VerifyAudience(audience, true) {
		return nil, errors.New("票据与访问平面不匹配")
	}
	return claims, nil
}

// verifyMachineTicket 校验服务间 S2S 机器票。只认 Authorization: Bearer。
// 密钥与受众都按目标团队选（与网关注入机器票时同一套映射）：一把密钥对应一个团队，
// 网关签发给 team-X 的票只能用 team-X 的密钥验过，验过即视为内部调用。
func (p Proxy) verifyMachineTicket(r *http.Request) bool {
	secret := strings.TrimSpace(p.machineSecret(r.URL.Path))
	if secret == "" {
		return false
	}
	raw := bearerToken(r)
	if raw == "" {
		return false
	}
	_, _, err := s2s.Verify(secret, p.machineAudience(r.URL.Path), raw)
	return err == nil
}

// identityHeader 按票面类型选择下游可信身份头名。
func identityHeader(kind ticketKind) string {
	if kind == ticketAdmin {
		return HeaderAdminID
	}
	return HeaderUserID
}

// injectMachineTicket 强制处理顺序第 5 步：用网关机器身份签发一张 S2S 票覆盖
// Authorization，使只挂 JwtS2s 的下游业务验过机器票后采信网关注入的身份头。
// 签发放 core 模块（非 middleware、非 issue.go），符合换票不进中间件的约束。
// 机器密钥未配置时返回 nil（开发态降级，生产必须配置）。
func (p Proxy) injectMachineTicket(r *http.Request) error {
	if p.routeAuthKind(r.URL.Path) == ticketPublic {
		return nil
	}
	m := p.Config.Machine
	secret := p.machineSecret(r.URL.Path)
	if strings.TrimSpace(secret) == "" {
		return nil
	}
	token, err := s2s.Issue(secret, m.Issuer, m.ClientID, p.machineAudience(r.URL.Path), nil, s2s.S2STTL)
	if err != nil {
		return err
	}
	r.Header.Set("Authorization", bearerPrefix+token)
	r.Header.Del(s2s.HeaderClientID)
	r.Header.Set(s2s.HeaderClientID, m.ClientID)
	return nil
}

// machineAudience 按 team/surface 选机器票受众：/beaver/admin/... → admin:internal，
// /beaver-open/api/... → open:internal，主站 /beaver/api/... → core:api。
// 先查 team/surface 精确键，再回退 team，最后回退默认。
func (p Proxy) machineAudience(path string) string {
	m := p.Config.Machine
	aud := m.DefaultAudience
	team, surface := teamAndSurface(path)
	if configured, ok := m.Audiences[team+"/"+surface]; ok {
		aud = configured
	} else if configured, ok := m.Audiences[team]; ok {
		aud = configured
	}
	return aud
}

// machineSecret 按 team/surface 选机器票密钥，查表口径与 machineAudience 完全一致：
// 先查 team/surface 精确键，再回退 team，最后回退 DefaultSecret / 旧 Secret。
// 密钥与受众用同一套路径映射，保证「发给 team-X 的票」只用 team-X 的密钥签发与校验。
func (p Proxy) machineSecret(path string) string {
	m := p.Config.Machine
	team, surface := teamAndSurface(path)
	if configured, ok := m.Secrets[team+"/"+surface]; ok {
		return configured
	}
	if configured, ok := m.Secrets[team]; ok {
		return configured
	}
	if strings.TrimSpace(m.DefaultSecret) != "" {
		return m.DefaultSecret
	}
	return m.Secret
}

// teamAndSurface 取六段路径的前两段（team/surface），用于选择机器票受众。
// 主站后台管理面 /beaver/admin/... 与 C 端 /beaver/api/... 同 team 不同 audience。
func teamAndSurface(path string) (string, string) {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

// bearerToken 只取 Authorization: Bearer 的裸票，用于机器票判定（机器票不走 Token 头 / query）。
func bearerToken(r *http.Request) string {
	raw := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(raw, bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(raw, bearerPrefix))
}

// ticketFromRequest 三种携带口径：Authorization: Bearer / Token 头 / query token。
func ticketFromRequest(r *http.Request) string {
	if raw := strings.TrimSpace(r.Header.Get("Authorization")); strings.HasPrefix(raw, bearerPrefix) {
		return strings.TrimSpace(strings.TrimPrefix(raw, bearerPrefix))
	}
	if t := strings.TrimSpace(r.Header.Get("Token")); t != "" {
		return t
	}
	return strings.TrimSpace(r.URL.Query().Get("token"))
}

// writeAuthError 鉴权失败统一 JSON 体（与全局 BaseResponse 一致）。
func writeAuthError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	body, _ := json.Marshal(BaseResponse{Code: 1, Msg: msg})
	_, _ = w.Write(body)
}
