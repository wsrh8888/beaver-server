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

package ai

import (
	"context"
	"net/http"

	"beaver/commonconfig"
	aitype "beaver/types/remote/ai"
)

// GetAgent 取 Agent 编排配置（人设 / 模型 / 工具；渠道无关的唯一真相源）。
// 走对方标准机器视图路径，不经网关。
func GetAgentApi(ctx context.Context, header http.Header, body *aitype.GetAgentReq) (*aitype.GetAgentRes, error) {
	return Post[aitype.GetAgentRes](ctx, commonconfig.AiUrl+"/beaver-ai/api/ai/aiconfig/v1/get_agent", header, body)
}
