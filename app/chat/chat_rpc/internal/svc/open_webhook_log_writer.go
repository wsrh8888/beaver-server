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
	"beaver/app/open/open_rpc/types/open_rpc"
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
	status := int32(0)
	if success {
		status = 1
	}
	_, err := w.openRpc.SaveWebhookLog(ctx, &open_rpc.SaveWebhookLogReq{
		ConfigId:  configID,
		AppId:     appID,
		EventType: eventType,
		Status:    status,
	})
	if err != nil {
		beaverlog.New("open_webhook_log", ctx).Error(model.LogMsg{Text: "SaveWebhookLog RPC 失败", Data: map[string]interface{}{"configId": configID, "appId": appID, "eventType": eventType, "err": err.Error()}})
	}
}
