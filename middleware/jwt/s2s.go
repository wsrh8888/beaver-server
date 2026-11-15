// Package jwt 服务间（S2S）机器票的鉴权，一个目录装齐「判定 + 接线」。
//
// 目录按「一件事」建，文件按层次分（与 middleware/track 的 httpTrack/rpcTrack 同一约定）：
//   - s2s.go         机器票（OAuth2 client_credentials，HS256）的签发 / 验签，纯密码学；
//   - httpS2S.go     机器票的 HTTP 闸门（读 Authorization 头、写 401、重建 X-Client-Id）。
//
// 网关签出的终端用户身份断言（ES256）本包当前不实现：本仓各面只校验机器票，
// 没有接线。将来哪个面要验断言，在本目录加 assertion.go（按 cloud 侧同名的签发格式反解），
// 需要 HTTP 接线再加 httpAssertion.go，不必另开顶层目录。
//
// 也不提供任何「网关跳令牌」闸门：网关只做转发，不往上游注入任何凭证，
// 因此没有可校验的「这一跳从网关来」。各面的身份只由业务凭证（app 钥匙 / 短票 / 回调 token）
// 或机器票决定。
//
// 出站换票（拿 client_id/secret 去对端 /token 换 Bearer 并缓存）不在本包，
// 在 remote 包 —— 它是调用方的连接状态，不是鉴权判定。
package jwt

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v4"
)

const (
	// S2STTL 换发的 S2S access token 默认有效期。
	S2STTL = 30 * time.Minute
	// HeaderClientID 验签通过后重建的、不可伪造的调用方标识头（写入下游 context 前剥离外部传入的同名头）。
	HeaderClientID = "X-Client-Id"
)

// S2SClaims 是 S2S access token 的载荷。
type S2SClaims struct {
	ClientID string   `json:"cid"` // 调用方 client_id
	Scopes   []string `json:"scp"` // 授权域（预留，按域授权时启用）
	jwtlib.StandardClaims
}

// Issue 用 secret 给 clientID 签发 audience=aud 的 HS256 JWT access token。
//
// issuer 标识签发方（一般取本团队服务名）；scopes 暂为空，按域授权时再填。
func Issue(secret, issuer, clientID, aud string, scopes []string, ttl time.Duration) (string, error) {
	if ttl <= 0 {
		ttl = S2STTL
	}
	now := time.Now()
	claims := S2SClaims{
		ClientID: clientID,
		Scopes:   scopes,
		StandardClaims: jwtlib.StandardClaims{
			Issuer:    issuer,
			Subject:   clientID,
			Audience:  aud,
			ExpiresAt: now.Add(ttl).Unix(),
			IssuedAt:  now.Unix(),
			Id:        randomHex(16),
		},
	}
	signed, err := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// Verify 验 HS256 JWT：签名 + 受众 + 过期。返回调用方 clientID 与 scopes。
func Verify(secret, aud, raw string) (string, []string, error) {
	parsed, err := jwtlib.ParseWithClaims(raw, &S2SClaims{}, func(t *jwtlib.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期签名算法: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", nil, err
	}
	c, ok := parsed.Claims.(*S2SClaims)
	if !ok || !parsed.Valid {
		return "", nil, errors.New("无效的令牌")
	}
	if !c.VerifyAudience(aud, true) {
		return "", nil, errors.New("令牌受众不匹配")
	}
	return c.ClientID, c.Scopes, nil
}

type ctxKey struct{}
type scopesCtxKey struct{}

// WithPrincipal 把已验签的调用方 clientID 注入 context（不可伪造：外部传入的同名头已被剥离）。
func WithPrincipal(ctx context.Context, clientID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, clientID)
}

// PrincipalFromContext 取已验签的调用方 clientID（无则空串）。
func PrincipalFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}

// WithScopes 把验签后票据里携带的授权域注入 context，供按域授权时读取。
func WithScopes(ctx context.Context, scopes []string) context.Context {
	return context.WithValue(ctx, scopesCtxKey{}, scopes)
}

// ScopesFromContext 取已验签票据里的授权域（无则 nil）。
func ScopesFromContext(ctx context.Context) []string {
	if v, ok := ctx.Value(scopesCtxKey{}).([]string); ok {
		return v
	}
	return nil
}

// ClientEntry 被调方入站注册表里的一个调用方条目。
// 注册表 map 的 key 是团队 / 配置段名，ClientID 才是换票请求里的凭据标识，两者不必相同。
type ClientEntry struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	Enabled      bool
}

// AuthenticateClient 校验 /token 端点的调用方：按 ClientID 找到条目、Enabled 为真、
// client_secret 常量时间比对一致。通过时返回该调用方被授予的 Scopes，供签发写进 S2SClaims。
//
// fail-closed：空 clientID / 空 clientSecret / 注册表里空口令一律拒绝，避免 ${VAR} 未注入时
// 退化成「两个空串比对相等」而放行。
func AuthenticateClient(clients map[string]ClientEntry, clientID, clientSecret string) ([]string, bool) {
	if clientID == "" || clientSecret == "" {
		return nil, false
	}
	for _, c := range clients {
		if c.ClientID != clientID || !c.Enabled || c.ClientSecret == "" {
			continue
		}
		return c.Scopes, subtle.ConstantTimeCompare([]byte(clientSecret), []byte(c.ClientSecret)) == 1
	}
	return nil, false
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
