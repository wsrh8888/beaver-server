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

// PeerListConversationsLogic 分页查会话
//
// 薄转发：HTTP 出入参与同域 RPC 的 proto 类型字段名一致，用 convert 做一次
// JSON 往返即可，不在这里复制业务规则（规则在 message_rpc 里）。
type PeerListConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPeerListConversationsLogic 构造 logic。
func NewPeerListConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PeerListConversationsLogic {
	return &PeerListConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PeerListConversations 分页查会话
func (l *PeerListConversationsLogic) PeerListConversations(req *types.PeerListConversationsReq) (*types.PeerListConversationsRes, error) {
	rpcReq := &chat_rpc.ListConversationsReq{UserId: req.UserID, ConversationType: req.ConversationType, Page: req.Page, PageSize: req.PageSize}
	rpcRes, err := l.svcCtx.ChatRpc.ListConversations(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp2 := make([]types.PeerConversationDetailItem, 0, len(rpcRes.List))
	for _, it1 := range rpcRes.List {
		tmp2 = append(tmp2, types.PeerConversationDetailItem{ConversationID: it1.ConversationId, Type: it1.Type, LastMessage: it1.LastMessage, LastMessageTime: it1.LastMessageTime, MessageCount: it1.MessageCount, ParticipantIDs: it1.ParticipantIds})
	}
	return &types.PeerListConversationsRes{Total: rpcRes.Total, List: tmp2}, nil
}
