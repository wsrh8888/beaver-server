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

// PeerGetConversationMemberIdsLogic 查会话成员 ID 列表
//
// 薄转发：出入参与同域 RPC 的 proto 字段名一致，用 convert 做一次 JSON 往返，
// 业务规则在 message_rpc 里。
type PeerGetConversationMemberIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPeerGetConversationMemberIdsLogic 构造 logic。
func NewPeerGetConversationMemberIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PeerGetConversationMemberIdsLogic {
	return &PeerGetConversationMemberIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PeerGetConversationMemberIds 查会话成员 ID 列表
func (l *PeerGetConversationMemberIdsLogic) PeerGetConversationMemberIds(req *types.PeerGetConversationMemberIdsReq) (*types.PeerGetConversationMemberIdsRes, error) {
	rpcReq := &chat_rpc.GetConversationMemberIDsReq{ConversationId: req.ConversationID, ExcludeUserId: req.ExcludeUserID}
	rpcRes, err := l.svcCtx.ChatRpc.GetConversationMemberIDs(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.PeerGetConversationMemberIdsRes{UserIDs: rpcRes.UserIds}, nil
}
