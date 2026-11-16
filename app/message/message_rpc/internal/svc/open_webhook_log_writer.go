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
	"beaver/core/corewebhook"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// openWebhookLogWriter 记录 webhook 投递结果到本地日志。
type openWebhookLogWriter struct{}

func newOpenWebhookLogWriter() corewebhook.LogWriter {
	return &openWebhookLogWriter{}
}

func (w *openWebhookLogWriter) SaveWebhookLog(ctx context.Context, configID, oauthID, eventType string, success bool) {
	// NOTE: 开放平台未提供 webhook 投递日志接口，此处以本地日志记录，configID 即 OAuthID。
	status := "fail"
	if success {
		status = "ok"
	}
	beaverlog.New("open_webhook_log", ctx).Info(model.LogMsg{
		Text: "webhook 投递日志（本地记录）",
		Data: map[string]interface{}{"configId": configID, "oauthId": oauthID, "eventType": eventType, "status": status},
	})
}
