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

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadByCursorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按分类游标标记已读（高效批量，首版主路径）
func NewMarkReadByCursorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadByCursorLogic {
	return &MarkReadByCursorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkReadByCursorLogic) MarkReadByCursor(req *types.MarkReadByCursorReq) (resp *types.MarkReadByCursorRes, err error) {
	resp = &types.MarkReadByCursorRes{Affected: 0}

	if req.UserID == "" || req.Category == "" {
		return resp, errors.New("userId 和 category 不能为空")
	}
	if req.ToVersion <= 0 && req.ToEventID == "" {
		return resp, errors.New("需提供 toVersion 或 toEventId")
	}

	now := time.Now()

	// 直接更新通知为已读
	query := l.svcCtx.DB.WithContext(l.ctx).Model(&notification_models.NotificationInbox{}).
		Where("user_id = ? AND category = ?", req.UserID, req.Category)

	if req.ToVersion > 0 {
		query = query.Where("version <= ?", req.ToVersion)
	} else if req.ToEventID != "" {
		query = query.Where("event_id = ?", req.ToEventID)
	}

	update := map[string]interface{}{
		"is_read": true,
		"read_at": now,
	}

	result := query.Updates(update)
	if result.Error != nil {
		return nil, result.Error
	}
	resp.Affected = result.RowsAffected

	return resp, nil
}
