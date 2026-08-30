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
	"time"

	"beaver/app/chat/chat_models"
	"beaver/app/chat/chat_rpc/internal/svc"
	"beaver/app/chat/chat_rpc/types/chat_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncMessageMediasLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSyncMessageMediasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncMessageMediasLogic {
	return &GetSyncMessageMediasLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSyncMessageMediasLogic) GetSyncMessageMedias(in *chat_rpc.GetSyncMessageMediasReq) (*chat_rpc.GetSyncMessageMediasRes, error) {
	var records []chat_models.ChatMessageMedia

	query := l.svcCtx.DB.Model(&chat_models.ChatMessageMedia{}).Where("user_id = ?", in.UserId)
	if in.Since > 0 {
		query = query.Where("UNIX_TIMESTAMP(created_at) > ?", in.Since)
	}

	if err := query.Find(&records).Error; err != nil {
		l.Errorf("同步消息媒体状态失败: userId=%s, error=%v", in.UserId, err)
		return nil, err
	}

	messageIDs := make([]string, 0, len(records))
	for _, item := range records {
		messageIDs = append(messageIDs, item.MessageID)
	}

	return &chat_rpc.GetSyncMessageMediasRes{
		MessageIds:      messageIDs,
		ServerTimestamp: time.Now().Unix(),
	}, nil
}
