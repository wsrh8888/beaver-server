// Package jwt 服务间（S2S）机器票的 HTTP 中间件入口，全仓唯一验票闸门。
//
// 纯逻辑（签发 / 验签 / 调用方凭据校验 / Principal context）在 utils/s2s，
// 本包只负责 HTTP 接线：读 Authorization 头、写 401、重建 X-Client-Id。
//
// 网关签出的终端用户身份断言（ES256）本包当前不实现：本仓各面只校验机器票，
// 没有接线。将来哪个面要验断言，在本目录加 assertion.go（按 cloud 侧同名的签发格式反解），
// 需要 HTTP 接线再加 httpAssertion.go，不必另开顶层目录。
//
// 终端用户身份不由本包判定：外部票（用户 / 开发者 / 管理员）只在网关验，
// 验过后写 Beaver-User-Id 或 Beaver-Admin-Id 可信头，再注入一张机器票覆盖 Authorization。
// 业务面只挂 JwtS2s 验这张机器票，验过即采信网关注入的身份头。不存在「网关跳令牌」式闸门——
// 可校验的东西就是机器票本身。
//
// 出站换票（拿 client_id/secret 去对端 /token 换 Bearer 并缓存）不在本包，
// 在 remote 包 —— 它是调用方的连接状态，不是鉴权判定。
package jwt

import (
	"net/http"
	"strings"

	"beaver/utils/s2s"
	"github.com/zeromicro/go-zero/rest"
)

const bearerPrefix = "Bearer "

// NewJwtS2s 机器票验签闸门（OAuth2 client_credentials + Bearer JWT），全仓唯一验票中间件入口。
//
// 口径：
//   - 调用方先经 /{team}/api/{service}/machine/v1/token 用 clientId + clientSecret 换 Bearer；
//   - 本中间件只认 `Authorization: Bearer <token>`，用 S2S.Secret 验 HS256 签名 + 受众 + 过期；
//     不接受终端用户票。
//   - 验过才剥离外部传入的 X-Client-Id，按验签结论重建，把 Principal 注入 context。
//
// audience 由各服务空壳按自身口径传入（主仓业务面统一用 s2s.AudienceCore）。
// fail-closed：密钥没配一律 401（服务 main 在启动期已先行退出）；令牌缺失 / 过期 /
// 受众不符同样 401。
func NewJwtS2s(secret, audience string) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return s2sGate(secret, audience)(next)
	}
}

// s2sGate 中间件本体，只依赖 http.HandlerFunc，验签走 utils/s2s.Verify。
func s2sGate(secret, aud string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if secret == "" {
				http.Error(w, "S2S 密钥未配置，拒绝访问", http.StatusUnauthorized)
				return
			}
			raw := strings.TrimSpace(r.Header.Get("Authorization"))
			if !strings.HasPrefix(raw, bearerPrefix) {
				http.Error(w, "缺少 Bearer 令牌", http.StatusUnauthorized)
				return
			}
			clientID, _, err := s2s.Verify(secret, aud, strings.TrimPrefix(raw, bearerPrefix))
			if err != nil {
				http.Error(w, "调用方令牌校验失败: "+err.Error(), http.StatusUnauthorized)
				return
			}
			// 重建不可伪造的调用方标识：剥离外部传入的同名头，按验签结论写入。
			r.Header.Del(s2s.HeaderClientID)
			r.Header.Set(s2s.HeaderClientID, clientID)
			// 票据里的 Scopes 来自签发时的注册表，验签通过后才可信；注入 context 供按域授权读取。
			next(w, r.WithContext(s2s.WithPrincipal(r.Context(), clientID)))
		}
	}
}
