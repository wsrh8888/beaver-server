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

package beaverlog

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"beaver/common/const/mqwsconst"
	"beaver/common/middleware/ua"
	"beaver/common/traceid"
	"beaver/core/corerocketmq"
	"beaver/utils/beaverlog/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

type ctxKey int

const (
	ctxKeyUserID ctxKey = iota + 1
	ctxKeyReq           // 请求上下文公参（打平上报，不进 message）
)

var (
	sourceMu sync.RWMutex
	source   string

	metaMu sync.RWMutex

	ipOnce   sync.Once
	serverIP string

	mqMu     sync.RWMutex
	mqClient *corerocketmq.Client
)

// Init 显式指定服务标识，如 auth_api、chat_rpc。
func Init(serviceSource string) {
	sourceMu.Lock()
	source = serviceSource
	sourceMu.Unlock()
}

// InitFromConf 从服务配置推导 module（yaml Name）。
func InitFromConf(c service.ServiceConf) {
	Init(strings.TrimSpace(c.Name))
}

// SetRocketMQ 注入 RocketMQ，启用扁平 JSON 上报（topic beaver_logs）。
func SetRocketMQ(c *corerocketmq.Client) {
	mqMu.Lock()
	mqClient = c
	mqMu.Unlock()
}

// Attach 把请求侧公参写入 ctx（入口处调一次即可）。
// r 非空时自动带上 http_request_url；并确保 ctx 有 traceId。
func Attach(ctx context.Context, r *http.Request, userID string, extra any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if userID != "" {
		ctx = context.WithValue(ctx, ctxKeyUserID, userID)
	}

	m := map[string]any{}
	if prev, ok := ctx.Value(ctxKeyReq).(map[string]any); ok {
		for k, v := range prev {
			m[k] = v
		}
	}
	if r != nil {
		m["http_request_url"] = r.URL.Path
	}
	for k, v := range toStringMap(extra) {
		m[k] = v
	}
	if len(m) > 0 {
		ctx = context.WithValue(ctx, ctxKeyReq, m)
		if id := mapString(m, "device_id"); id != "" {
			ctx = context.WithValue(ctx, ua.KeyDeviceID, id)
		}
	}

	if traceid.FromContext(ctx) == "" {
		var serverTrace, clientUuid string
		if r != nil {
			serverTrace = traceid.ServerTraceFromHeaders(r.Header)
			clientUuid = traceid.ClientUuidFromHeaders(r.Header)
		}
		ctx = traceid.WithContext(ctx, traceid.Ensure(serverTrace), clientUuid)
	}
	return ctx
}

// Middleware HTTP 入口挂公参。顺序：TraceMiddleware → ua.Middleware → beaverlog.Middleware。
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := Attach(r.Context(), r, r.Header.Get("Beaver-User-Id"), map[string]any{
			"device_id":          r.Header.Get("Beaver-Device-Id"),
			"client_remote_addr": r.RemoteAddr,
			"user_agent":         r.Header.Get("User-Agent"),
		})
		next(w, r.WithContext(ctx))
	}
}

func currentSource() string {
	sourceMu.RLock()
	defer sourceMu.RUnlock()
	return source
}

func currentMQ() *corerocketmq.Client {
	mqMu.RLock()
	defer mqMu.RUnlock()
	return mqClient
}

// Logger module 为服务内子模块（仅本地 logx）；上报 module 用服务名。
type Logger struct {
	module string
	ctx    context.Context
}

// New 创建日志器。ctx 可选：beaverlog.New("module") 或 beaverlog.New("module", ctx)
func New(module string, ctx ...context.Context) *Logger {
	l := &Logger{module: module}
	if len(ctx) > 0 {
		l.ctx = ctx[0]
	}
	return l
}

func (l *Logger) Info(msg model.LogMsg) {
	l.send("info", msg)
}

func (l *Logger) Warn(msg model.LogMsg) {
	l.send("warn", msg)
}

func (l *Logger) Error(msg model.LogMsg) {
	l.send("error", msg)
}

// logPayload 顶层扁平字段；message 仅业务 text/data。
type logPayload struct {
	Module           string `json:"module"`
	Level            string `json:"level"`
	TraceID          string `json:"traceId"`
	UserID           string `json:"user_id"`
	DeviceID         string `json:"device_id"`
	ClientRemoteAddr string `json:"client_remote_addr"`
	HTTPRequestURL   string `json:"http_request_url"`
	UserAgent        string `json:"user_agent"`
	Message          string `json:"message"`
	Timestamp        int64  `json:"timestamp"`
}

type messageBody struct {
	Text string      `json:"text"`
	Data interface{} `json:"data,omitempty"`
}

func (l *Logger) send(level string, msg model.LogMsg) {
	req := reqFromCtx(l.ctx)
	userID := userIDFromCtx(l.ctx)
	traceID := traceid.FromContext(l.ctx)
	deviceID := mapString(req, "device_id")

	fields := []logx.LogField{
		logx.Field("module", currentSource()),
	}
	if l.module != "" {
		fields = append(fields, logx.Field("scene", l.module))
	}
	if userID != "" {
		fields = append(fields, logx.Field("user_id", userID))
	}
	if traceID != "" {
		fields = append(fields, logx.Field("traceId", traceID))
	}
	if deviceID != "" {
		fields = append(fields, logx.Field("device_id", deviceID))
	}
	if msg.Data != nil {
		fields = append(fields, logx.Field("data", msg.Data))
	}

	ctx := l.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	lg := logx.WithContext(ctx).WithCallerSkip(1)

	switch level {
	case "warn":
		lg.Sloww(msg.Text, fields...)
	case "error":
		lg.Errorw(msg.Text, fields...)
	default:
		lg.Infow(msg.Text, fields...)
	}

	l.publishMQ(level, msg, userID, traceID, req)
}

func (l *Logger) publishMQ(level string, msg model.LogMsg, userID, traceID string, req map[string]any) {
	mq := currentMQ()
	if mq == nil {
		return
	}

	messageBytes, err := json.Marshal(messageBody{Text: msg.Text, Data: msg.Data})
	if err != nil {
		messageBytes = []byte(msg.Text)
	}

	payload := logPayload{
		Module:           currentSource(),
		Level:            level,
		TraceID:          traceID,
		UserID:           userID,
		DeviceID:         mapString(req, "device_id"),
		ClientRemoteAddr: mapString(req, "client_remote_addr"),
		HTTPRequestURL:   mapString(req, "http_request_url"),
		UserAgent:        mapString(req, "user_agent"),
		Message:          string(messageBytes),
		Timestamp:        time.Now().UnixMilli(),
	}

	go func() {
		if err := mq.SendRawJSON(context.Background(), mqwsconst.MqTopicClientLog, payload); err != nil {
			logx.Debugf("beaverlog mq publish failed: %v", err)
		}
	}()
}

func userIDFromCtx(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v, _ := ctx.Value(ctxKeyUserID).(string)
	return v
}

func reqFromCtx(ctx context.Context) map[string]any {
	m := map[string]any{}
	if ctx == nil {
		return m
	}
	if v, ok := ctx.Value(ctxKeyReq).(map[string]any); ok {
		for k, val := range v {
			m[k] = val
		}
	}
	if _, ok := m["device_id"]; !ok {
		if v, ok := ctx.Value(ua.KeyDeviceID).(string); ok && v != "" {
			m["device_id"] = v
		}
	}
	return m
}

func toStringMap(v any) map[string]any {
	if v == nil {
		return nil
	}
	if m, ok := v.(map[string]any); ok {
		out := make(map[string]any, len(m))
		for k, val := range m {
			if s, ok := val.(string); ok && s == "" {
				continue
			}
			out[k] = val
		}
		if len(out) == 0 {
			return nil
		}
		return out
	}
	b, err := json.Marshal(v)
	if err != nil || len(b) == 0 || string(b) == "null" || string(b) == "{}" {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil || len(m) == 0 {
		return nil
	}
	return m
}

func mapString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
