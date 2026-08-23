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

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadByEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按事件ID标记单个通知已读
func NewMarkReadByEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadByEventLogic {
	return &MarkReadByEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkReadByEventLogic) MarkReadByEvent(req *types.MarkReadByEventReq) (resp *types.MarkReadByEventRes, err error) {
	userId := req.UserID
	eventId := req.EventID

	// 更新指定通知为已读
	result := l.svcCtx.DB.Model(&notification_models.NotificationInbox{}).
		Where("user_id = ? AND event_id = ? AND is_deleted = ?", userId, eventId, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": time.Now(),
		})

	if result.Error != nil {
		l.Logger.Errorf("标记单个通知已读失败: %v", result.Error)
		return nil, result.Error
	}

	resp = &types.MarkReadByEventRes{}

	return resp, nil
}
