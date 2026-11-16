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

package messagelogic

import (
	"context"
	"time"

	message_models "beaver/app/message/message_models/chat"
	"beaver/app/message/message_rpc/internal/svc"
	"beaver/app/message/message_rpc/pb"

	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncDeletedMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetSyncDeletedMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncDeletedMessagesLogic {
	return &GetSyncDeletedMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("GetSyncDeletedMessagesLogic", ctx),
	}
}

func (l *GetSyncDeletedMessagesLogic) GetSyncDeletedMessages(in *pb.GetSyncDeletedMessagesReq) (*pb.GetSyncDeletedMessagesRes, error) {
	var deletes []message_models.ChatUserDelete

	query := l.svcCtx.DB.Model(&message_models.ChatUserDelete{}).Where("user_id = ?", in.UserId)
	if in.Since > 0 {
		query = query.Where("UNIX_TIMESTAMP(created_at) > ?", in.Since)
	}
	if err := query.Find(&deletes).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "同步已删除消息失败", Data: map[string]any{"userId": in.UserId, "err": err.Error()}})
		return nil, err
	}

	msgIDs := make([]string, 0, len(deletes))
	for _, d := range deletes {
		msgIDs = append(msgIDs, d.MessageID)
	}
	return &pb.GetSyncDeletedMessagesRes{
		MessageIds:      msgIDs,
		ServerTimestamp: time.Now().Unix(),
	}, nil
}
