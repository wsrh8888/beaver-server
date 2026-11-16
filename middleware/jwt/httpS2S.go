package jwt

import (
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

// NewS2S 机器票验签闸门（OAuth2 client_credentials + Bearer JWT）。
//
// 口径（OAuth2 client_credentials + Bearer JWT）：
//   - 调用方先经 /beaver/api/auth/token/v1/issue 用 clientId + clientSecret 换 Bearer；
//   - 本中间件只认 `Authorization: Bearer <token>`，用 S2S.Secret 验 HS256 签名 + 受众
//   - 过期；不接受终端用户票。
//   - 验过才剥离外部传入的 X-Client-Id，按验签结论重建，把 Principal 注入 context。
//
// fail-closed：密钥没配一律 401（服务 main 在启动期已先行退出）；令牌缺失 / 过期 /
// 受众不符同样 401。
func NewS2S(secret, aud string) func(http.HandlerFunc) http.HandlerFunc {
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
			clientID, _, err := Verify(secret, aud, strings.TrimPrefix(raw, bearerPrefix))
			if err != nil {
				http.Error(w, "调用方令牌校验失败: "+err.Error(), http.StatusUnauthorized)
				return
			}
			// 重建不可伪造的调用方标识：剥离外部传入的同名头，按验签结论写入。
			r.Header.Del(HeaderClientID)
			r.Header.Set(HeaderClientID, clientID)
			// 票据里的 Scopes 来自签发时的注册表，验签通过后才可信；注入 context 供按域授权读取。
			ctx := WithPrincipal(r.Context(), clientID)
			next(w, r.WithContext(ctx))
		}
	}
}
