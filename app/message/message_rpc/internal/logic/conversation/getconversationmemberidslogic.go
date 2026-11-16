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

package conversationlogic

import (
	"context"

	message_models "beaver/app/message/message_models/chat"
	"beaver/app/message/message_rpc/internal/svc"
	"beaver/app/message/message_rpc/pb"

	beaverlog "beaver/utils/beaverlog"
)

type GetConversationMemberIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetConversationMemberIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationMemberIDsLogic {
	return &GetConversationMemberIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetConversationMemberIDsLogic", ctx),
	}
}

func (l *GetConversationMemberIDsLogic) GetConversationMemberIDs(in *pb.GetConversationMemberIDsReq) (*pb.GetConversationMemberIDsRes, error) {
	var rows []message_models.ChatUserConversation
	query := l.svcCtx.DB.Model(&message_models.ChatUserConversation{}).Where("conversation_id = ?", in.ConversationId)
	if in.ExcludeUserId != "" {
		query = query.Where("user_id <> ?", in.ExcludeUserId)
	}
	if err := query.Find(&rows).Error; err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(rows))
	for _, row := range rows {
		userIDs = append(userIDs, row.UserID)
	}
	return &pb.GetConversationMemberIDsRes{UserIds: userIDs}, nil
}
