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

package agent

import (
	"context"
	"net/http"

	"beaver/commonconfig"
	agenttype "beaver/types/remote/agent"
)

// SendMessage 发一句话给 Agent（fire-and-forget：只取 code，不等终稿）。
func SendMessageApi(ctx context.Context, header http.Header, body *agenttype.SendMessageReq) (*agenttype.SendMessageRes, error) {
	return Post[agenttype.SendMessageRes](ctx, commonconfig.AgentUrl+"/beaver-ai/api/agent/chat/v1/send_message", header, body)
}
