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

package traceid

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

const (
	// HeaderRequestID 服务端链路 ID（Gateway 生成，与客户端无关）
	HeaderRequestID = "X-Request-Id"
	// HeaderClientUuid 客户端自己的请求标识（仅透传/记录，不参与服务端 trace）
	HeaderClientUuid = "Uuid"
	// MetadataKey gRPC metadata 键（服务端 trace）
	MetadataKey = "x-request-id"
)

type ctxKey struct{}
type clientCtxKey struct{}

// New 生成服务端 traceId。
func New() string {
	return uuid.NewString()
}

// Ensure 空则生成（仅用于服务端 ID）。
func Ensure(id string) string {
	if id != "" {
		return id
	}
	return New()
}

// ClientUuidFromHeaders 只读客户端 Uuid，绝不当作服务端 trace。
func ClientUuidFromHeaders(h http.Header) string {
	return h.Get(HeaderClientUuid)
}

// ServerTraceFromHeaders 只读服务端 X-Request-Id（不看 Uuid）。
func ServerTraceFromHeaders(h http.Header) string {
	return h.Get(HeaderRequestID)
}

// AttachServerTrace 写入服务端 trace 到请求头/响应头；不改动客户端 Uuid。
func AttachServerTrace(w http.ResponseWriter, r *http.Request, serverTraceID string) {
	r.Header.Set(HeaderRequestID, serverTraceID)
	if w != nil {
		w.Header().Set(HeaderRequestID, serverTraceID)
	}
}

// WithContext 注入服务端 trace；可选附带客户端 uuid（仅日志字段，互不覆盖）。
func WithContext(ctx context.Context, serverTraceID, clientUuid string) context.Context {
	if serverTraceID != "" {
		ctx = context.WithValue(ctx, ctxKey{}, serverTraceID)
		ctx = logx.ContextWithFields(ctx, logx.Field("trace", serverTraceID))
	}
	if clientUuid != "" {
		ctx = context.WithValue(ctx, clientCtxKey{}, clientUuid)
		ctx = logx.ContextWithFields(ctx, logx.Field("clientUuid", clientUuid))
	}
	return ctx
}

// FromContext 读取服务端 traceId。
func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(ctxKey{}).(string); ok {
		return v
	}
	return ""
}

// ClientUuidFromContext 读取客户端 Uuid。
func ClientUuidFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(clientCtxKey{}).(string); ok {
		return v
	}
	return ""
}

// FromIncomingMD 从入站 gRPC metadata 取服务端 traceId。
func FromIncomingMD(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals := md.Get(MetadataKey)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}
