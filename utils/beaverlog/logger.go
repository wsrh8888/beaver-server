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
	"fmt"
	"strings"
	"sync"
	"time"

	"beaver/utils/beaverlog/model"

	"go.opentelemetry.io/otel/attribute"
	logapi "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var (
	sourceMu     sync.RWMutex
	source       string
	otlpAddr     string // 远程日志(OTLP)采集器地址；空字符串表示关闭远程日志（默认）
	remoteOnce   sync.Once
	remoteLogger logapi.Logger // nil 表示未启用远程日志（默认关闭）
)

// Init 显式指定服务标识，如 auth_api、chat_rpc。
// 新服务推荐直接用 InitFromConf，避免服务名两处维护。
func Init(serviceSource string) {
	sourceMu.Lock()
	source = serviceSource
	sourceMu.Unlock()
}

// InitFromConf 从服务配置推导 source，自动带上，杜绝漏调 Init 导致 source 为空。
// 直接取 Name（yaml 顶层 Name，已统一为目录约定名，如 emoji_api / emoji_rpc）。
// 用法：beaverlog.InitFromConf(c.RestConf) 或 beaverlog.InitFromConf(c.RpcConf)
func InitFromConf(c service.ServiceConf) {
	Init(strings.TrimSpace(c.Name))
}

func currentSource() string {
	sourceMu.RLock()
	defer sourceMu.RUnlock()
	return source
}

// SetOtlpAddr 设置远程日志(OTLP)采集器地址，如 127.0.0.1:4317。
// 空字符串表示关闭远程日志（默认）。须在首次打日志前调用（通常在服务启动时、InitFromConf 之后）。
func SetOtlpAddr(addr string) {
	otlpAddr = strings.TrimSpace(addr)
}

// Logger module 为服务内子模块；ctx 可选
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

func (l *Logger) send(level string, msg model.LogMsg) {
	fields := []logx.LogField{
		logx.Field("source", currentSource()),
	}
	if l.module != "" {
		fields = append(fields, logx.Field("module", l.module))
	}
	if msg.Data != nil {
		fields = append(fields, logx.Field("data", msg.Data))
	}

	ctx := l.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	// WithCallerSkip(1) 抵消 send 这一层包装，否则 caller 恒指向本文件的 Infow/Sloww/Errorw 行，
	// 线上无法定位到真正打日志的业务代码。
	lg := logx.WithContext(ctx).WithCallerSkip(1)

	switch level {
	case "warn":
		lg.Sloww(msg.Text, fields...)
	case "error":
		lg.Errorw(msg.Text, fields...)
	default:
		lg.Infow(msg.Text, fields...)
	}

	// 远程推送（OTel OTLP -> 采集器 -> OpenSearch）。
	// 仅当启动时通过 beaverlog.SetOtlpAddr 设置了采集器地址才启用，默认关闭，
	// 因此只给需要观测的服务在 yaml 里配置 OtlpAddr 即可，不影响其他服务。
	remoteOnce.Do(initRemote)
	if remoteLogger != nil {
		l.emitRemote(level, msg)
	}
}

// initRemote 惰性初始化 OTel OTLP 日志上报。仅 SetOtlpAddr 设置了采集器地址才开启。
func initRemote() {
	addr := otlpAddr
	if addr == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exp, err := otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(addr),
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithTimeout(3*time.Second),
	)
	if err != nil {
		// 远程日志初始化失败不影响本地日志
		return
	}
	res, err := resource.New(ctx, resource.WithAttributes(
		attribute.String("service.name", currentSource()),
	))
	if err != nil || res == nil {
		res = resource.Default()
	}
	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
	)
	remoteLogger = provider.Logger("beaverlog")
}

// emitRemote 将一条日志以 OTLP 格式异步上报（BatchProcessor 内部已做批处理与后台发送）。
func (l *Logger) emitRemote(level string, msg model.LogMsg) {
	var sev logapi.Severity
	switch level {
	case "error":
		sev = logapi.SeverityError1
	case "warn":
		sev = logapi.SeverityWarn1
	default:
		sev = logapi.SeverityInfo1
	}

	var rec logapi.Record
	rec.SetTimestamp(time.Now())
	rec.SetSeverity(sev)

	// body 优先取客户端日志正文（Data.message）；透传场景下 msg.Text 恒为 "客户端日志上报"，
	// 真正的正文在 Data 内，故以 Data.message 覆盖，避免 Discover 里每条都是同一句占位文本。
	bodyText := msg.Text
	if m, ok := msg.Data.(map[string]any); ok {
		if msgVal, ok := m["message"]; ok {
			if s, ok := msgVal.(string); ok && s != "" {
				bodyText = s
			}
		}
	}
	rec.SetBody(attribute.StringValue(bodyText))

	// 平铺 Data 为独立 attribute，使 Dashboards 可按 platform / deviceId / level / userId 等字段过滤。
	// source / module 由 beaverlog 注入；若 Data 自带同名键则以其为准，避免重复 key。
	// Data 非 map（如结构体）时退化为整体 data 字符串，保证服务端日志原有行为兼容。
	attrs := make([]attribute.KeyValue, 0, 8)
	if m, ok := msg.Data.(map[string]any); ok && m != nil {
		if _, ok := m["source"]; !ok {
			attrs = append(attrs, attribute.String("source", currentSource()))
		}
		if _, ok := m["module"]; !ok {
			attrs = append(attrs, attribute.String("module", l.module))
		}
		for k, v := range m {
			if k == "source" || k == "module" {
				continue // 已按上述规则处理，杜绝重复 key
			}
			attrs = append(attrs, toLogAttr(k, v))
		}
	} else {
		attrs = append(attrs, attribute.String("source", currentSource()))
		attrs = append(attrs, attribute.String("module", l.module))
		if msg.Data != nil {
			if b, err := json.Marshal(msg.Data); err == nil {
				attrs = append(attrs, attribute.String("data", string(b)))
			}
		}
	}
	rec.AddAttributes(attrs...)

	ctx := l.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	remoteLogger.Emit(ctx, rec)
}

// toLogAttr 将任意值安全转换为 OTLP attribute。
// JSON 反序列化后数字为 float64、对象为 map、数组为 slice，这里做最小类型推导；
// 不支持的结构退化为 JSON 字符串，确保字段始终可被索引与检索。
func toLogAttr(key string, val any) attribute.KeyValue {
	switch v := val.(type) {
	case string:
		return attribute.String(key, v)
	case bool:
		return attribute.Bool(key, v)
	case float64:
		return attribute.Float64(key, v)
	case float32:
		return attribute.Float64(key, float64(v))
	case int:
		return attribute.Int64(key, int64(v))
	case int64:
		return attribute.Int64(key, v)
	case int32:
		return attribute.Int64(key, int64(v))
	case nil:
		return attribute.String(key, "")
	default:
		if b, err := json.Marshal(v); err == nil {
			return attribute.String(key, string(b))
		}
		return attribute.String(key, fmt.Sprintf("%v", v))
	}
}
