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

package sse

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Emitter 向客户端推一帧 SSE 事件。
type Emitter func(event string, payload any) error

// WriteHeaders 设置 SSE 响应头（须在首次写 body 前调用）。
func WriteHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}

// Upgrade 校验 Flusher 并写好 SSE 头，返回可复用的 Emitter。
func Upgrade(w http.ResponseWriter) (Emitter, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("streaming unsupported")
	}
	WriteHeaders(w)
	return func(event string, payload any) error {
		return Write(w, flusher, event, payload)
	}, nil
}

// Write 写一帧：event + data(JSON) + flush。
func Write(w http.ResponseWriter, flusher http.Flusher, event string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data); err != nil {
		return err
	}
	flusher.Flush()
	return nil
}
