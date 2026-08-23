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

	"beaver/app/notification/notification_api/internal/svc"
	"beaver/app/notification/notification_api/internal/types"
	"beaver/app/notification/notification_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetReadCursorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按分类拉取通知已读游标
func NewGetReadCursorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetReadCursorsLogic {
	return &GetReadCursorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetReadCursorsLogic) GetReadCursors(req *types.GetReadCursorsReq) (resp *types.GetReadCursorsRes, err error) {
	resp = &types.GetReadCursorsRes{Cursors: []types.ReadCursorItem{}}

	if req.UserID == "" {
		return resp, nil
	}

	var rows []notification_models.NotificationRead
	query := l.svcCtx.DB.WithContext(l.ctx).
		Where("user_id = ?", req.UserID)

	if len(req.Categories) > 0 {
		query = query.Where("category IN ?", req.Categories)
	}

	if err = query.Find(&rows).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, row := range rows {
		lastReadAt := int64(0)
		if row.LastReadAt != nil {
			lastReadAt = row.LastReadAt.UnixMilli()
		}

		resp.Cursors = append(resp.Cursors, types.ReadCursorItem{
			Category:   row.Category,
			Version:    row.Version,
			LastReadAt: lastReadAt,
		})
	}

	return resp, nil
}
