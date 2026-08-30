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

type GetInboxByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按ID拉取收件箱明细
func NewGetInboxByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInboxByIdsLogic {
	return &GetInboxByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetInboxByIdsLogic) GetInboxByIds(req *types.GetInboxByIdsReq) (resp *types.GetInboxByIdsRes, err error) {
	resp = &types.GetInboxByIdsRes{Inbox: []types.InboxItem{}}

	if len(req.EventIDs) == 0 || req.UserID == "" {
		return resp, nil
	}

	var rows []notification_models.NotificationInbox
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Where("user_id = ? AND event_id IN ?", req.UserID, req.EventIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		var readAt int64
		if row.ReadAt != nil {
			readAt = time.Time(*row.ReadAt).UnixMilli()
		}

		resp.Inbox = append(resp.Inbox, types.InboxItem{
			EventID:   row.EventID,
			EventType: row.EventType,
			Category:  row.Category,
			Version:   row.Version,
			IsRead:    row.IsRead,
			ReadAt:    readAt,
			Status:    int32(row.Status),
			IsDeleted: row.IsDeleted,
			Silent:    row.Silent,
			CreatedAt: time.Time(row.CreatedAt).UnixMilli(),
			UpdatedAt: time.Time(row.UpdatedAt).UnixMilli(),
		})
	}

	return resp, nil
}
