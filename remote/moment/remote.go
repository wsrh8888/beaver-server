package moment

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"beaver/utils/ajax"
)

// 本包唯一的出站通道：换票 + 带票发包，其余不碰（不重试、不聚合、不缓存业务结果）。
// 各业务线的对端、凭据、路径各不相同，所以这份 plumbing 每个业务线各存一份，不抽公共包。
//
// 换票请求在 token.go 的 GetToken，域名由业务方法在自己的调用行写全
// commonconfig.SocialUrl + 路径。

const (
	// tokenTTL 对端没给 expires_in 时的兜底有效期。
	tokenTTL = 30 * time.Minute
	// tokenSkew 提前量：快过期就重换，避免请求发出时刚好失效。
	tokenSkew = 30 * time.Second
)

// Post 发一次 POST 的 JSON 请求，返回值是信封 result 解成的 Res。
func Post[Res, Req any](ctx context.Context, endpoint string, header http.Header, body Req) (*Res, error) {
	h, err := client.authHeader(ctx, header)
	if err != nil {
		return nil, fmt.Errorf("%s：换取出站令牌失败: %w", endpoint, err)
	}
	return ajax.Post[Res](ctx, endpoint, "", h, body)
}

// Get 发一次 GET，q 为 nil 表示不带查询串。
func Get[Res any](ctx context.Context, endpoint string, header http.Header, q url.Values) (*Res, error) {
	h, err := client.authHeader(ctx, header)
	if err != nil {
		return nil, fmt.Errorf("%s：换取出站令牌失败: %w", endpoint, err)
	}
	if len(q) > 0 {
		endpoint += "?" + q.Encode()
	}
	return ajax.Get[Res](ctx, endpoint, "", h)
}

// Delete 发一次 DELETE：go-zero 的 delete 路由同样从 body 取参数。
func Delete[Res, Req any](ctx context.Context, endpoint string, header http.Header, body Req) (*Res, error) {
	h, err := client.authHeader(ctx, header)
	if err != nil {
		return nil, fmt.Errorf("%s：换取出站令牌失败: %w", endpoint, err)
	}
	return ajax.Delete[Res](ctx, endpoint, "", h, body)
}

// authHeader 把调用方给的头复制一份，再补上 Authorization。
func (c *Client) authHeader(ctx context.Context, header http.Header) (http.Header, error) {
	bearer, err := c.bearerToken(ctx)
	if err != nil {
		return nil, err
	}
	h := http.Header{}
	for k, vs := range header {
		for _, v := range vs {
			h.Set(k, v)
		}
	}
	h.Set("Authorization", bearer)
	return h, nil
}

// bearerToken 取本次请求要带的 Authorization 值，必要时先去换票。
// 缓存到过期前 30 秒；锁住整个换票过程，让同时到达的一批请求只打一次 /token。
func (c *Client) bearerToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.bearer != "" && time.Now().Before(c.expiresAt.Add(-tokenSkew)) {
		return c.bearer, nil
	}
	tok, exp, err := c.exchange(ctx)
	if err != nil {
		return "", err
	}
	c.bearer, c.expiresAt = tok, exp
	return tok, nil
}

// exchange 换密钥：凭据交给 token.go 的 GetToken 换一次票，这里只算过期时间。
func (c *Client) exchange(ctx context.Context) (string, time.Time, error) {
	token, expiresIn, err := GetTokenApi(ctx, c.clientID, c.clientSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	ttl := tokenTTL
	if expiresIn > 0 {
		ttl = time.Duration(expiresIn) * time.Second
	}
	return "Bearer " + token, time.Now().Add(ttl), nil
}
