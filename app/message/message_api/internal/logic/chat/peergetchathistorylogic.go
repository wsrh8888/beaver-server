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
	"errors"
	"slices"
	"strings"

	accttypes "beaver/app/account/account_rpc/pb"
	"beaver/app/message/message_api/internal/svc"
	"beaver/app/message/message_api/internal/types"
	message_models "beaver/app/message/message_models/chat"
	chat_rpc "beaver/app/message/message_rpc/pb"
	"beaver/common/list_query"
	"beaver/common/models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

// PeerGetChatHistoryLogic 按会话取可读对话流，供摘要 / Agent 之类需要"人话"
// 的消费方使用。
//
// 与用户面 ChatHistory 的差别只在身份来源：S2S 面没有网关注入的
// Beaver-User-Id，身份由请求体 UserID 携带，所以必须在此校验它确为会话成员，
// 否则一张内部票就能读任意会话。除此之外口径与用户面一致：同样过滤掉该用户
// 主动删除的消息，同样只取正常状态的消息。
type PeerGetChatHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewPeerGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PeerGetChatHistoryLogic {
	return &PeerGetChatHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("peer_getchathistory", ctx),
	}
}

const peerHistoryDefaultLimit = 20
const peerHistoryMaxLimit = 100

func (l *PeerGetChatHistoryLogic) PeerGetChatHistory(req *types.PeerChatHistoryReq) (*types.PeerChatHistoryRes, error) {
	uid := strings.TrimSpace(req.UserID)
	conversationID := strings.TrimSpace(req.ConversationID)
	if uid == "" || conversationID == "" {
		return nil, errors.New("userId and conversationId required")
	}

	// 越权闸门：身份来自请求体，必须确认它真是这个会话的成员。
	members, err := l.svcCtx.ChatRpc.GetConversationMemberIDs(l.ctx, &chat_rpc.GetConversationMemberIDsReq{
		ConversationId: conversationID,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查会话成员失败",
			Data: map[string]any{"userId": uid, "conversationId": conversationID, "err": err.Error()},
		})
		return nil, err
	}
	if !slices.Contains(members.GetUserIds(), uid) {
		return nil, errors.New("user is not a member of this conversation")
	}

	limit := req.Limit
	if limit <= 0 {
		limit = peerHistoryDefaultLimit
	}
	if limit > peerHistoryMaxLimit {
		limit = peerHistoryMaxLimit
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}

	// 过滤当前用户主动删除的消息（与用户面 ChatHistory 同口径）。
	subQuery := l.svcCtx.DB.Model(&message_models.ChatUserDelete{}).
		Select("message_id").
		Where("user_id = ?", uid)
	where := l.svcCtx.DB.Where("conversation_id = ? AND message_id NOT IN (?) AND status = 1", conversationID, subQuery)
	if req.BeforeSeq > 0 {
		where = where.Where("seq < ?", req.BeforeSeq)
	}
	rows, total, err := list_query.ListQuery(l.svcCtx.DB, message_models.ChatMessage{}, list_query.Option{
		PageInfo: models.PageInfo{Page: page, Limit: limit, Sort: "seq desc"},
		Where:    where,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询聊天历史失败",
			Data: map[string]any{"userId": uid, "conversationId": conversationID, "err": err.Error()},
		})
		return nil, err
	}

	// 批量取昵称：消费方要的是"谁说的"，不是裸 userId。
	userIDs := make([]string, 0, len(rows))
	seen := map[string]bool{}
	nickMap := map[string]string{}
	for _, row := range rows {
		if row.SendUserID != nil && *row.SendUserID != "" && !seen[*row.SendUserID] {
			seen[*row.SendUserID] = true
			userIDs = append(userIDs, *row.SendUserID)
		}
	}
	if len(userIDs) > 0 {
		if resp, e := l.svcCtx.UserRpc.UserListInfo(l.ctx, &accttypes.UserListInfoReq{UserIdList: userIDs}); e == nil {
			for id, info := range resp.UserInfo {
				nickMap[id] = info.NickName
			}
		}
	}

	list := make([]types.PeerChatHistoryItem, 0, len(rows))
	for _, row := range rows {
		senderID := ""
		if row.SendUserID != nil {
			senderID = *row.SendUserID
		}
		nick := nickMap[senderID]
		if senderID == "" {
			nick = "通知"
		} else if nick == "" {
			nick = senderID
		}
		list = append(list, types.PeerChatHistoryItem{
			MessageID:  row.MessageID,
			SenderID:   senderID,
			SenderName: nick,
			Content:    row.MsgPreviewMethod(),
			MsgType:    uint32(row.MsgType),
			Seq:        row.Seq,
			CreatedAt:  row.CreatedAt.String(),
		})
	}

	return &types.PeerChatHistoryRes{
		Total:   total,
		HasMore: int64(page*limit) < total,
		List:    list,
	}, nil
}
