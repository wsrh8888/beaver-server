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

	chattypes "beaver/app/message/message_rpc/pb"
	relationtypes "beaver/app/relation/relation_rpc/pb"
	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	type_struct "beaver/app/ws/ws_api/types"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
	"beaver/utils/beaverlog/model"
	"beaver/utils/conversation"
)

func HandleTypingSend(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	req *types.WsReq,
	_ *http.Request,
	_ *ws_conn.Client,
	bodyRaw json.RawMessage,
) error {
	var body type_struct.BodyTyping
	if err := json.Unmarshal(bodyRaw, &body); err != nil {
		return fmt.Errorf("typing 消息格式错误: %w", err)
	}
	if body.ConversationID == "" {
		return fmt.Errorf("conversationId 不能为空")
	}

	peerIDs, err := getTypingPeerIDs(ctx, svcCtx, req.UserID, body.ConversationID)
	if err != nil {
		return err
	}

	for _, peerID := range peerIDs {
		payload := map[string]interface{}{
			"command":  wsCommandConst.CHAT_MESSAGE,
			"type":     wsTypeConst.TypingReceive,
			"senderId": req.UserID,
			"targetId": peerID,
			"body": map[string]interface{}{
				"conversationId": body.ConversationID,
				"userId":         req.UserID,
				"isTyping":       body.IsTyping,
			},
			"conversationId": body.ConversationID,
		}
		if err := svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logger.Error(model.LogMsg{
				Text: "推送typing通知失败",
				Data: map[string]any{"sender": req.UserID, "target": peerID, "err": err.Error()},
			})
		}
	}
	return nil
}

func getTypingPeerIDs(ctx context.Context, svcCtx *svc.ServiceContext, currentUserID, conversationID string) ([]string, error) {
	convType, userIDs := conversation.ParseConversationWithType(conversationID)
	if convType == 1 {
		var peers []string
		for _, uid := range userIDs {
			if uid != currentUserID {
				peers = append(peers, uid)
			}
		}
		if len(peers) == 0 {
			return nil, fmt.Errorf("无效的私聊会话")
		}
		return peers, nil
	}

	groupID := conversation.GetTargetIDByConversation(conversationID, currentUserID)
	membersRes, err := svcCtx.GroupRpc.GetGroupMembers(ctx, &relationtypes.GetGroupMembersReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, fmt.Errorf("查询群成员失败: %w", err)
	}

	peers := make([]string, 0, len(membersRes.Members))
	for _, m := range membersRes.Members {
		if m.UserID != currentUserID {
			peers = append(peers, m.UserID)
		}
	}
	if len(peers) == 0 {
		// 群成员为空时的兜底：会话成员属 message 库，走 message_rpc，不再直连
		memberRes, err := svcCtx.ChatRpc.GetConversationMemberIDs(ctx, &chattypes.GetConversationMemberIDsReq{
			ConversationId: conversationID,
			ExcludeUserId:  currentUserID,
		})
		if err != nil {
			return nil, err
		}
		peers = append(peers, memberRes.UserIds...)
	}
	return peers, nil
}
