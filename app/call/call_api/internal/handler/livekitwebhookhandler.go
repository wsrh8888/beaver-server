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

package handler

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"beaver/app/call/call_api/internal/logic"
	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/common/response"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/webhook"
)

func LiveKitWebhookHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	authProvider := auth.NewSimpleKeyProvider(
		svcCtx.Config.LiveKit.ApiKey,
		svcCtx.Config.LiveKit.ApiSecret,
	)

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		if _, err := webhook.ReceiveWebhookEvent(r, authProvider); err != nil {
			response.Response(r, w, nil, errors.New("webhook签名验证失败"))
			return
		}

		req := types.LiveKitWebhookReq{Body: body}
		l := logic.NewLiveKitWebhookLogic(r.Context(), svcCtx)
		resp, err := l.LiveKitWebhook(&req)
		response.Response(r, w, resp, err)
	}
}
