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
	"time"

	"beaver/app/chat/chat_models"
	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncDeletedMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取所有需要更新的已删除消息ID列表
func NewGetSyncDeletedMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncDeletedMessagesLogic {
	return &GetSyncDeletedMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncDeletedMessagesLogic) GetSyncDeletedMessages(req *types.GetSyncDeletedMessagesReq) (resp *types.GetSyncDeletedMessagesRes, err error) {
	var deletes []chat_models.ChatUserDelete

	// 1. 查询该用户在 Since 之后的删除记录
	query := l.svcCtx.DB.Model(&chat_models.ChatUserDelete{}).Where("user_id = ?", req.UserID)
	if req.Since > 0 {
		// 统一使用 Unix 时间戳进行同步比对
		query = query.Where("UNIX_TIMESTAMP(created_at) > ?", req.Since)
	}

	err = query.Find(&deletes).Error
	if err != nil {
		l.Logger.Errorf("同步已删除消息失败: userId=%s, error=%v", req.UserID, err)
		return nil, errors.New("同步失败")
	}

	// 2. 提取消息ID列表
	msgIDs := make([]string, 0, len(deletes))
	for _, d := range deletes {
		msgIDs = append(msgIDs, d.MessageID)
	}

	return &types.GetSyncDeletedMessagesRes{
		MessageIDs:      msgIDs,
		ServerTimestamp: time.Now().Unix(),
	}, nil
}
