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

package logger

import (
	"sync"

	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	sourceMu sync.RWMutex
	source   string
)

// Init 每个服务启动时调用一次，source 标识当前微服务（如 auth_api、chat_api）
func Init(serviceSource string) {
	sourceMu.Lock()
	source = serviceSource
	sourceMu.Unlock()
}

func currentSource() string {
	sourceMu.RLock()
	defer sourceMu.RUnlock()
	return source
}

// Logger module 为服务内子模块，source 由 Init 统一设置
type Logger struct {
	module string
}

func New(module string) *Logger {
	return &Logger{module: module}
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

	switch level {
	case "warn":
		logx.Sloww(msg.Text, fields...)
	case "error":
		logx.Errorw(msg.Text, fields...)
	default:
		logx.Infow(msg.Text, fields...)
	}
}
