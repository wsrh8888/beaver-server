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

package chat_message

import (
	"context"
	"fmt"

	"beaver/app/open/open_api/internal/svc"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsTypeConst"
)

// Handle 机器人上行 CHAT_MESSAGE 业务帧的分派入口。
//
// 与用户侧 ws_api 的 chat_message.Handle 结构一致：按 content.Data.Type 分派。
// 机器人当前只有一种上行业务类型（流式增量）—— 机器人的「回复」不走 WS，
// 统一走 OAPI（POST /api/open/robot/v1/send_message），这样出站逻辑只有一份。
func Handle(ctx context.Context, svcCtx *svc.ServiceContext, robotID string, content type_struct.WsContent) error {
	switch content.Data.Type {
	case wsTypeConst.ChatMessageStreamSend:
		return HandleStreamSend(ctx, svcCtx, robotID, content)
	default:
		// 不再静默丢弃：带上具体 type 打日志，否则龙虾侧调试时会以为帧发出去了
		return fmt.Errorf("不支持的机器人上行消息类型: %s", content.Data.Type)
	}
}
