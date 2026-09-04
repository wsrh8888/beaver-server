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

package grpcMiddleware

import (
	"context"
	"time"

	"beaver/common/middleware/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// RequestLogInterceptor gRPC 请求日志拦截器
func RequestLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	code := 0
	if err != nil {
		if st, ok := status.FromError(err); ok {
			code = int(st.Code())
		} else {
			code = 13 // Internal
		}
	}
	utils.LogRequest(ctx, "gRPC", info.FullMethod, req, resp, err, code, startTime)
	return resp, err
}
