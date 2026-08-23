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
	"errors"

	logic "beaver/app/auth/auth_api/internal/logic/auth_public"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/middleware/ua"
	"beaver/common/response"
	"beaver/utils/device"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PhoneLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PhoneLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		if err := validateLoginDevice(r, req.DeviceID); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewPhoneLoginLogic(r.Context(), svcCtx)
		resp, err := l.PhoneLogin(&req)
		response.Response(r, w, resp, err)
	}
}

func validateLoginDevice(r *http.Request, deviceID string) error {
	preciseType := ua.DeviceType(r.Context())
	if preciseType == "" || preciseType == device.DeviceUnknown {
		return errors.New("不支持的设备类型")
	}
	if deviceID == "" {
		return errors.New("无法识别的物理设备，请联系管理员")
	}
	return nil
}
