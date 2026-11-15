// Package ajax 出站 HTTP 的统一封装（内网服务间 / 外部系统调用）。
//
// 对应前端的 utils/request：只管「怎么把请求发出去、怎么判成败」，不管「发到哪个
// 路径、字段是什么意思」——后者在 remote/<模块> 的调用封装里按接口分开写。
//
// 统一三件事，都是被真实故障逼出来的（对齐 cloud-server/utils/ajax）：
//   - 禁代理：etcd 注册的是内网 IP（192.168.x / 10.x），Go 的 loopback 跳过规则不生效，
//     默认 Transport 会把内网请求塞进宿主 HTTP_PROXY，表现为「上游明明活着却 503」；
//   - 共享连接池：每个客户端各建一个 Transport，高频路径上会反复建连（TIME_WAIT 打满端口）；
//   - 信封判定：go-zero 的 {code,msg,result}，HTTP 200 + code!=0 也是失败，
//     判定成败必须看 code，不能只看状态码。
//
// 这里不管凭据：请求头由调用方给全，本包不换 token、不认 client_id/secret。跨团队
// 换票属于对端契约，落在 remote/<模块> 的调用封装里（那边换来 Bearer 再交给这里）。
//
// 不做重试：每个下游的幂等语义不同，重试归属调用方。
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

// StatusError 下游返回了非成功响应：HTTP 非 2xx，或 HTTP 200 + code!=0。
//
// 调用方用 errors.As 取 Status 自行判定（如 401 = 凭据不被对端认可）。
type StatusError struct {
	Name    string // 下游标识（取自请求路径）
	Status  int    // HTTP 状态码
	Code    int    // 信封业务码（未解析出信封时为 0）
	Msg     string // 信封 msg（下游没给原因时为空）
	NonJSON bool   // 响应体不是信封 JSON（文件下载、裸 401 页等）
}

func (e *StatusError) Error() string {
	switch {
	case e.Status == http.StatusUnauthorized:
		return e.Name + "：下游调用方凭据校验失败（401），检查对端配置的 client_id / 公钥与本侧签发凭据是否成对"
	case e.Status < 200 || e.Status >= 300:
		if e.Msg != "" {
			return fmt.Sprintf("%s：下游接口返回 %d：%s", e.Name, e.Status, e.Msg)
		}
		if e.NonJSON {
			return fmt.Sprintf("%s：下游接口返回 %d（响应非信封 JSON）", e.Name, e.Status)
		}
		return fmt.Sprintf("%s：下游接口返回 %d", e.Name, e.Status)
	default:
		return fmt.Sprintf("%s：下游返回失败：%s", e.Name, FailReason(e.Msg))
	}
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

// decode 按 go-zero 信封 {code,msg,result} 判成败并解出 result。Post / Get 共用。
func decode[Res any](path string, raw []byte, status int) (*Res, error) {
	var env struct {
		Code   int             `json:"code"`
		Msg    string          `json:"msg"`
		Result json.RawMessage `json:"result"`
	}
	parsed := json.Unmarshal(raw, &env) == nil

	if status < 200 || status >= 300 {
		se := &StatusError{Name: path, Status: status, NonJSON: !parsed}
		if parsed {
			se.Code, se.Msg = env.Code, env.Msg
		}
		return nil, se
	}
	if !parsed {
		return nil, fmt.Errorf("%s：解析下游响应失败（非信封 JSON）", path)
	}
	if env.Code != 0 {
		return nil, &StatusError{Name: path, Status: status, Code: env.Code, Msg: env.Msg}
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
		se := &StatusError{Name: url, Status: status, NonJSON: json.Unmarshal(raw, &oauthErr) != nil}
		se.Msg = oauthErr.ErrorDescription
		if se.Msg == "" {
			se.Msg = oauthErr.Error
		}
		return nil, se
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

// FailReason 下游没给文案时给个兜底，绝不把空 message 透给客户。
func FailReason(msg string) string {
	if strings.TrimSpace(msg) == "" {
		return "下游未返回原因，请稍后重试或联系运营"
	}
	return msg
}
