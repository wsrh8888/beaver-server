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

package utils

import (
	"time"

	"beaver/utils/logger"
	"beaver/utils/logger/model"
)

var accessLog = logger.New("http")

// LogRequest 记录 HTTP/gRPC 请求日志
func LogRequest(method, path string, req, resp interface{}, err error, startTime time.Time) {
	duration := time.Since(startTime)
	data := map[string]interface{}{
		"method":   method,
		"path":     path,
		"duration": duration.String(),
	}
	if req != nil {
		data["req"] = req
	}
	if resp != nil {
		data["resp"] = resp
	}

	if err != nil {
		data["err"] = err.Error()
		accessLog.Error(model.LogMsg{
			Text: "请求失败",
			Data: data,
		})
		return
	}

	accessLog.Info(model.LogMsg{
		Text: "请求成功",
		Data: data,
	})
}
