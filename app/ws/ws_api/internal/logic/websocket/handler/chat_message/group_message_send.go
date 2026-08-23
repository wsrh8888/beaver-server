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
	"encoding/json"
	"fmt"
	"net/http"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)

var groupMsgLog = logger.New("group_msg_send")

func HandleGroupMessageSend(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *types.WsReq,
	r *http.Request,
	client *ws_conn.Client,
	messageId string,
	bodyRaw json.RawMessage,
) error {
	var body type_struct.BodySendMsg
	if err := json.Unmarshal(bodyRaw, &body); err != nil {
		groupMsgLog.Error(model.LogMsg{
			Text: "群聊消息解析失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return fmt.Errorf("消息格式错误: %w", err)
	}

	rpcMsg, err := convertToRpcMsg(body.Msg)
	if err != nil {
		groupMsgLog.Error(model.LogMsg{
			Text: "群聊消息格式转换失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return fmt.Errorf("消息内容错误: %w", err)
	}

	_, err = svcCtx.ChatRpc.SendMsg(ctx, &chat_rpc.SendMsgReq{
		UserId:         req.UserID,
		ConversationId: body.ConversationID,
		MessageId:      messageId,
		Msg:            rpcMsg,
	})
	if err != nil {
		groupMsgLog.Error(model.LogMsg{
			Text: "群聊消息发送失败",
			Data: map[string]interface{}{
				"userId":         req.UserID,
				"conversationId": body.ConversationID,
				"err":            err.Error(),
			},
		})
		return err
	}

	return nil
}
