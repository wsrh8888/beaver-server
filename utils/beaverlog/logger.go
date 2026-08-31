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
	"strings"
	"sync"

	"beaver/utils/beaverlog/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var (
	sourceMu sync.RWMutex
	source   string
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
}
