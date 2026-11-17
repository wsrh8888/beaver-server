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
package logic

import (
	"context"

	"beaver/app/message/message_api/internal/svc"
	"beaver/app/message/message_api/internal/types"
	chat_rpc "beaver/app/message/message_rpc/pb"
)

// PeerListChatMessagesLogic 分页查消息
//
// 薄转发：HTTP 出入参与同域 RPC 的 proto 类型字段名一致，用 convert 做一次
// JSON 往返即可，不在这里复制业务规则（规则在 message_rpc 里）。
type PeerListChatMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPeerListChatMessagesLogic 构造 logic。
func NewPeerListChatMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PeerListChatMessagesLogic {
	return &PeerListChatMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PeerListChatMessages 分页查消息
func (l *PeerListChatMessagesLogic) PeerListChatMessages(req *types.PeerListChatMessagesReq) (*types.PeerListChatMessagesRes, error) {
	rpcReq := &chat_rpc.ListChatMessagesReq{
		Page: req.Page, PageSize: req.PageSize, ConversationId: req.ConversationID,
		SendUserId: req.SendUserID, MsgType: req.MsgType, StartTime: req.StartTime,
		EndTime: req.EndTime, MessageId: req.MessageID, WithContent: req.WithContent,
		Status: req.Status, Order: req.Order,
	}
	rpcRes, err := l.svcCtx.ChatRpc.ListChatMessages(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	res := &types.PeerListChatMessagesRes{Total: rpcRes.Total}
	for _, it := range rpcRes.GetList() {
		res.List = append(res.List, types.PeerChatMessageItem{
			MessageID: it.MessageId, ConversationID: it.ConversationId,
			ConversationType: it.ConversationType, SendUserID: it.SendUserId,
			MsgType: it.MsgType, MsgPreview: it.MsgPreview, MsgContent: it.MsgContent,
			Status: it.Status, CreatedAt: it.CreatedAt, UpdatedAt: it.UpdatedAt,
		})
	}
	return res, nil
}
