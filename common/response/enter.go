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

package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

// Response allows for a custom message to be passed in case of success
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error, successMsg ...string) {
	msg := ""
	if len(successMsg) > 0 {
		msg = successMsg[0]
	}
	if err == nil {
		r := &Body{
			Code:   0,
			Msg:    msg,
			Result: resp,
		}
		httpx.WriteJson(w, http.StatusOK, r)
		return
	}
	httpx.WriteJson(w, http.StatusOK, &Body{
		Code:   1,
		Msg:    err.Error(),
		Result: nil,
	})
}
