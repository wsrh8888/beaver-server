// Package ajax 出站 HTTP 统一封装：禁代理（内网直连）、共享连接池、按 {code,msg,result} 信封判成败。
// 只管「怎么发、怎么判成败」；路径与字段含义在 remote/<模块> 里按接口写。
// 凭据由调用方给全（header），重试归调用方（下游幂等语义各异）。
package ajax

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"beaver/common/traceid"
)

// Transport 内网出站共享连接池。内网 S2S 并发量级低一档，池子按 32/主机足够。
var Transport = &http.Transport{
	// 内网一律直连：服务发现拿到的地址都在可信内网，不经代理是更正确的一侧。
	Proxy: nil,
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          256,
	MaxIdleConnsPerHost:   32,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// callTimeout 单次调用超时。内网 S2S 都是同步短请求，10s 足够，不做成可配项。
const callTimeout = 10 * time.Second

// HTTPClient 只取传输能力（共享池 + 禁代理 + 超时），不解信封、不管头。
func HTTPClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &http.Client{
		Timeout:   timeout,
		Transport: Transport,
	}
}

// envelope 统一响应信封：成功 {code,msg,result}（code=0），失败 {code,msg}（result 缺失）。
type envelope struct {
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Result json.RawMessage `json:"result"`
}

// StatusError 下游返回失败：HTTP 非 2xx，或 HTTP 200 + code!=0。
type StatusError struct {
	Status int    // HTTP 状态码
	Code   int    // 信封业务码
	Msg    string // 信封 msg
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("下游返回失败：status=%d code=%d msg=%s", e.Status, e.Code, e.Msg)
}

// Post 发一次 POST 的 JSON 请求，返回值是信封 result 解成的 Res。
//
// header 是调用方给好的完整请求头（含 Authorization 等凭据头），可为 nil。换 token 不在这里：
// 那是 remote/<模块> 的职责，它换好 Bearer 塞进 header 再调这里。
func Post[Res, Req any](ctx context.Context, baseURL, path string, header http.Header, body Req) (*Res, error) {
	return Do[Res](ctx, http.MethodPost, baseURL, path, header, body)
}

// Delete 发一次 DELETE 的 JSON 请求。go-zero 的 delete 路由同样从 body 取参数。
func Delete[Res, Req any](ctx context.Context, baseURL, path string, header http.Header, body Req) (*Res, error) {
	return Do[Res](ctx, http.MethodDelete, baseURL, path, header, body)
}

// Do 发一次带 JSON body 的请求（POST / DELETE），返回值是信封 result 解成的 Res。
func Do[Res, Req any](ctx context.Context, method, baseURL, path string, header http.Header, body Req) (*Res, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("%s：下游接口地址未配置", path)
	}

	var payload []byte
	if !isNil(body) {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%s：序列化请求失败: %w", path, err)
		}
		payload = b
	}

	req, err := NewJSONRequest(ctx, method, strings.TrimRight(baseURL, "/")+path, payload)
	if err != nil {
		return nil, fmt.Errorf("%s：构造请求失败: %w", path, err)
	}
	applyHeader(req, header)

	raw, status, err := send(req)
	if err != nil {
		return nil, fmt.Errorf("%s：%w", path, err)
	}
	return decode[Res](path, raw, status)
}

// Get 发一次 GET 请求（无请求体），响应按信封解成 Res。
func Get[Res any](ctx context.Context, baseURL, path string, header http.Header) (*Res, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("%s：下游接口地址未配置", path)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, strings.TrimRight(baseURL, "/")+path, nil)
	if err != nil {
		return nil, fmt.Errorf("%s：构造请求失败: %w", path, err)
	}
	traceid.SetOutgoingHeader(ctx, req.Header)
	applyHeader(req, header)

	raw, status, err := send(req)
	if err != nil {
		return nil, fmt.Errorf("%s：%w", path, err)
	}
	return decode[Res](path, raw, status)
}

// decode 按 {code,msg,result} 信封判成败并解出 result。Post / Get / Delete 共用。
func decode[Res any](path string, raw []byte, status int) (*Res, error) {
	var env envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("%s：解析下游响应失败: %w", path, err)
	}
	if status < 200 || status >= 300 || env.Code != 0 {
		return nil, &StatusError{Status: status, Code: env.Code, Msg: env.Msg}
	}

	out := new(Res)
	if len(env.Result) > 0 && string(env.Result) != "null" {
		if err := json.Unmarshal(env.Result, out); err != nil {
			return nil, fmt.Errorf("%s：解析下游 result 失败: %w", path, err)
		}
	}
	return out, nil
}

// applyHeader 把调用方给好的头写进请求，跳过空值（空值表示「这一项没配」，不是「清空」）。
func applyHeader(req *http.Request, header http.Header) {
	for k, vs := range header {
		for _, v := range vs {
			if strings.TrimSpace(v) != "" {
				req.Header.Set(k, v)
			}
		}
	}
}

// send 发一次请求并读回响应体（限 64KB）。只管传输，不判成败、不解信封。
func send(req *http.Request) ([]byte, int, error) {
	resp, err := HTTPClient(callTimeout).Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("调用下游接口失败: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("读取下游响应失败: %w", err)
	}
	return raw, resp.StatusCode, nil
}

// PostForm 发一次表单请求（application/x-www-form-urlencoded），响应按裸 JSON 解成 Res。
//
// 只给「不套信封」的端点用：OAuth2 换票（access_token / expires_in）返回的就是裸 JSON，
// 没有 go-zero 的 {code,msg,result}。业务接口一律走 Post。
func PostForm[Res any](ctx context.Context, url string, header http.Header, form url.Values) (*Res, error) {
	if strings.TrimSpace(url) == "" {
		return nil, fmt.Errorf("表单请求：下游接口地址未配置")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("表单请求：构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	traceid.SetOutgoingHeader(ctx, req.Header)
	applyHeader(req, header)

	raw, status, err := send(req)
	if err != nil {
		return nil, fmt.Errorf("表单请求：%w", err)
	}
	if status < 200 || status >= 300 {
		var oauthErr struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		_ = json.Unmarshal(raw, &oauthErr)
		msg := oauthErr.ErrorDescription
		if msg == "" {
			msg = oauthErr.Error
		}
		return nil, &StatusError{Status: status, Msg: msg}
	}
	out := new(Res)
	if err := json.Unmarshal(raw, out); err != nil {
		return nil, fmt.Errorf("表单请求：解析下游响应失败: %w", err)
	}
	return out, nil
}

// isNil 判断请求体是否为空。泛型下的 typed nil（如 (*Req)(nil)）不能直接和 nil 比，
// 不判就会把字符串 "null" 当 body 发给下游，报错点会跑到下游去。
func isNil(v any) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
		return rv.IsNil()
	}
	return false
}

// NewJSONRequest 构造 JSON 请求：统一 Content-Type 与 trace 头（X-Request-Id）。
func NewJSONRequest(ctx context.Context, method, url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	traceid.SetOutgoingHeader(ctx, req.Header)
	return req, nil
}
