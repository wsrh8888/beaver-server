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
	"net/http"

	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsTypeConst"
)

func Handle(ctx context.Context, svcCtx *svc.ServiceContext, req *types.WsReq, r *http.Request, client *ws_conn.Client, content type_struct.WsContent) error {
	switch content.Data.Type {
	case wsTypeConst.GroupMessageSend:
		return HandleGroupMessageSend(ctx, svcCtx, req, r, client, content.MessageID, content.Data.Body)
	case wsTypeConst.PrivateMessageSend:
		return HandlePrivateMessageSend(ctx, svcCtx, req, r, client, content.MessageID, content.Data.Body)
	case wsTypeConst.TypingSend:
		return HandleTypingSend(ctx, svcCtx, req, r, client, content.Data.Body)
	default:
		fmt.Println("未支持的消息类型:", content.Data.Type)
		return fmt.Errorf("unsupported message type: %s", content.Data.Type)
	}
}
