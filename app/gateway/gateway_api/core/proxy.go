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
	"beaver/app/gateway/gateway_api/types"
	"beaver/common/etcd"
	"beaver/common/traceid"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"beaver/utils/jwts"
	utils "beaver/utils/list"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// apiGwLogger 网关 API 代理包级日志器（用于无 ctx 的辅助函数）
var apiGwLogger = beaverlog.New("gateway_api_proxy")

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func writeErrorResponse(res http.ResponseWriter, msg string, statusCode int, uuid string) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	response := BaseResponse{Code: 1, Msg: msg}
	byteData, _ := json.Marshal(response)
	res.Write(byteData)
	apiGwLogger.Info(model.LogMsg{Text: "响应输出", Data: map[string]interface{}{"uuid": uuid, "body": string(byteData)}})

}

type Proxy struct {
	Config types.Config
}

func (p Proxy) auth(req *http.Request) (ok bool, errMsg string) {
	path := req.URL.Path

	// 1. 公开接口，无需鉴权
	if utils.InListByRegex(p.Config.PublicList, path) {
		return true, ""
	}

	// 2. 自定义鉴权，透传到下游服务 middleware 处理
	if utils.InListByRegex(p.Config.CustomAuthList, path) {
		return true, ""
	}

	// 3. open oauth_secret：Gateway 校验 App-Id / App-Secret 请求头
	if strings.HasPrefix(path, "/api/open/oauth_secret/") {
		return p.oauthSecretAuth(req)
	}

	// 4. open_api：Gateway 默认不鉴权，各接口 logic 自行校验；仅 Beaver JWT 路由例外
	if isOpenApiPassThrough(path) {
		return true, ""
	}

	// 5. 统一 JWT 鉴权
	if !p.jwtAuth(req) {
		return false, "网关鉴权失败"
	}
	return true, ""
}

func (p Proxy) oauthSecretAuth(req *http.Request) (bool, string) {
	appID := req.Header.Get("App-Id")
	if appID == "" {
		return false, "缺少 App-Id 请求头"
	}

	appSecret := req.Header.Get("App-Secret")
	if appSecret == "" {
		return false, "缺少 App-Secret 请求头"
	}
	return true, ""
}

var openApiJwtRoutes = []string{
	`/api/open/oauth/v1/h5_authcode`,
	`/api/open/oauth/v1/qrcode_scan`,
	`/api/open/oauth/v1/qrcode_confirm`,
	`/api/open/oauth/v1/qrcode_cancel`,
}

func isOpenApiPassThrough(path string) bool {
	if !strings.HasPrefix(path, "/api/open/") {
		return false
	}
	if utils.InListByRegex(openApiJwtRoutes, path) {
		return false
	}
	return true
}

// jwtAuth JWT认证（普通用户）
func (p Proxy) jwtAuth(req *http.Request) bool {
	// 获取token
	token := getToken(req)
	if token == "" {
		apiGwLogger.Error(model.LogMsg{Text: "token为空"})
		return false
	}

	// 直接解析JWT（避免HTTP调用）
	claims, err := jwts.ParseToken(token, p.Config.Auth.AccessSecret)
	if err != nil {
		apiGwLogger.Error(model.LogMsg{Text: "JWT解析失败", Data: map[string]interface{}{"err": err.Error()}})
		return false
	}

	// 设置用户ID和设备ID到请求头
	req.Header.Set("Beaver-User-Id", claims.UserID)

	// 从请求头获取设备ID
	deviceId := req.Header.Get("deviceId")
	if deviceId != "" {
		req.Header.Set("Beaver-Device-Id", deviceId)
	}
	version := req.Header.Get("version")
	if version != "" {
		req.Header.Set("Version", version)
	}

	apiGwLogger.Info(model.LogMsg{Text: "JWT验证成功", Data: map[string]interface{}{"userId": claims.UserID, "deviceId": deviceId}})

	return true
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	_, _ = r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func bodyForLog(raw []byte) any {
	if len(raw) == 0 {
		return nil
	}
	const maxLogBody = 8 << 10 // 8KB，避免大文件/大包打爆日志
	truncated := false
	if len(raw) > maxLogBody {
		raw = raw[:maxLogBody]
		truncated = true
	}
	var obj any
	if err := json.Unmarshal(raw, &obj); err == nil {
		if truncated {
			return map[string]any{"data": obj, "truncated": true}
		}
		return obj
	}
	s := string(raw)
	if truncated {
		return s + "...(truncated)"
	}
	return s
}

func (p Proxy) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	clientUuid := traceid.ClientUuidFromHeaders(req.Header)
	serverTrace := ensureServerTraceId(res, req)
	ctx := traceid.WithContext(req.Context(), serverTrace, clientUuid)

	lg := beaverlog.New("proxy", ctx)
	rec := &statusRecorder{ResponseWriter: res, status: http.StatusOK}
	clientIP := getClientIP(req)
	method := req.Method
	path := req.URL.Path

	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		_ = req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
		req.ContentLength = int64(len(reqBody))
	}

	finish := func(text string, level string, extra map[string]any) {
		data := map[string]any{
			"method":   method,
			"path":     path,
			"clientIp": clientIP,
			"status":   rec.status,
			"duration": time.Since(start).String(),
			"req":      bodyForLog(reqBody),
			"resp":     bodyForLog(rec.body.Bytes()),
		}
		if q := req.URL.RawQuery; q != "" {
			data["query"] = q
		}
		for k, v := range extra {
			data[k] = v
		}
		msg := model.LogMsg{Text: text, Data: data}
		switch level {
		case "warn":
			lg.Warn(msg)
		case "error":
			lg.Error(msg)
		default:
			lg.Info(msg)
		}
	}

	// 限流检查
	if p.Config.Limit.Enable {
		if !p.rateLimitCheck(clientIP) {
			writeErrorResponse(rec, "请求频率过高", http.StatusTooManyRequests, serverTrace)
			finish("请求频率过高", "warn", nil)
			return
		}
	}

	token := getToken(req)
	req.Header.Set("Token", token)
	ok, authErrMsg := p.auth(req)
	if !ok {
		writeErrorResponse(rec, authErrMsg, http.StatusUnauthorized, serverTrace)
		finish("网关鉴权失败", "warn", map[string]any{"msg": authErrMsg})
		return
	}

	// 匹配路由
	regex, _ := regexp.Compile(`/api/(.*?)/`)

	addrList := regex.FindStringSubmatch(req.URL.Path)
	if len(addrList) != 2 {
		writeErrorResponse(rec, "请求不匹配", http.StatusBadRequest, serverTrace)
		finish("请求不匹配", "warn", nil)
		return
	}

	service := addrList[1]

	// 增加重试机制
	var addr string
	for i := 0; i < 3; i++ {
		addr = etcd.GetServiceAddr(p.Config.Etcd, service+"_api")
		if addr != "" {
			break
		}
		lg.Error(model.LogMsg{Text: "获取服务地址失败", Data: map[string]interface{}{"attempt": i + 1, "service": service + "_api"}})
		time.Sleep(100 * time.Millisecond)
	}

	if addr == "" {
		writeErrorResponse(rec, "服务暂时不可用", http.StatusServiceUnavailable, serverTrace)
		finish("服务不可用", "error", map[string]any{"service": service + "_api"})
		lg.Error(model.LogMsg{Text: "未匹配到服务", Data: map[string]interface{}{"service": service + "_api"}})
		return
	}

	remote, _ := url.Parse(fmt.Sprintf("http://%s", addr))
	reverseProxy := httputil.NewSingleHostReverseProxy(remote)

	// 修改默认的 Director 函数以保留 User-Agent
	originalDirector := reverseProxy.Director
	reverseProxy.Director = func(req *http.Request) {
		originalDirector(req)
		if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}
		if clientIP := getClientIP(req); clientIP != "" {
			req.Header.Set("ClientIp", clientIP)
		}
		// 只透传服务端 X-Request-Id；客户端 Uuid 原样保留，互不覆盖
		if id := req.Header.Get(traceid.HeaderRequestID); id != "" {
			req.Header.Set(traceid.HeaderRequestID, id)
		}
	}

	reverseProxy.ServeHTTP(rec, req)
	finish("网关请求完成", "info", map[string]any{"service": service + "_api"})
}

func getToken(req *http.Request) string {
	token := req.Header.Get("Token")
	if token == "" {
		token = req.URL.Query().Get("token")
	}
	return token
}

// ensureServerTraceId Gateway 入口始终生成服务端 trace，与客户端 Uuid 无关。
func ensureServerTraceId(res http.ResponseWriter, req *http.Request) string {
	id := traceid.New()
	traceid.AttachServerTrace(res, req, id)
	return id
}

// getClientIP 获取客户端真实IP
func getClientIP(req *http.Request) string {
	// 优先从 X-Forwarded-For 获取
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// 其次从 X-Real-IP 获取
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// 最后从 RemoteAddr 获取
	return req.RemoteAddr
}

// rateLimitCheck 限流检查
var globalRateLimiter *RateLimiter

func (p Proxy) rateLimitCheck(clientIP string) bool {
	if globalRateLimiter == nil {
		globalRateLimiter = NewRateLimiter(p.Config.Limit.Rate, p.Config.Limit.Burst)
	}
	return globalRateLimiter.Allow(clientIP)
}
