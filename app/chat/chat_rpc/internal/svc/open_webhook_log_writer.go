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

package svc

import (
	"context"

	"beaver/app/open/open_rpc/open"
	"beaver/core/corewebhook"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type openRpcWebhookLogWriter struct {
	openRpc open.Open
}

func newOpenRpcWebhookLogWriter(openRpc open.Open) corewebhook.LogWriter {
	return &openRpcWebhookLogWriter{openRpc: openRpc}
}

func (w *openRpcWebhookLogWriter) SaveWebhookLog(ctx context.Context, configID, appID, eventType string, success bool) {
	// NOTE: open_rpc 在「开放平台 Digital/Bot 模型与门户能力重构」中移除了 SaveWebhookLog RPC，
	// 此处暂以本地日志替代，待 open 侧补回日志接口后再恢复 RPC 调用。
	status := "fail"
	if success {
		status = "ok"
	}
	beaverlog.New("open_webhook_log", ctx).Info(model.LogMsg{
		Text: "webhook 日志（open_rpc 暂无 SaveWebhookLog 接口，本地记录）",
		Data: map[string]interface{}{"configId": configID, "appId": appID, "eventType": eventType, "status": status},
	})
}
