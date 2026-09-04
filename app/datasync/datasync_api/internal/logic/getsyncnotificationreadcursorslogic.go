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

	"beaver/app/datasync/datasync_api/internal/svc"
	"beaver/app/datasync/datasync_api/internal/types"
	"beaver/app/notification/notification_rpc/types/notification_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetSyncNotificationReadCursorsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取通知已读游标版本摘要
func NewGetSyncNotificationReadCursorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncNotificationReadCursorsLogic {
	return &GetSyncNotificationReadCursorsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_sync_notification_read_cursors", ctx),
	}
}

func (l *GetSyncNotificationReadCursorsLogic) GetSyncNotificationReadCursors(req *types.GetSyncNotificationReadCursorsReq) (resp *types.GetSyncNotificationReadCursorsRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcResp, err := l.svcCtx.NotificationRpc.GetReadCursorVersions(l.ctx, &notification_rpc.GetReadCursorVersionsReq{
		UserId:       req.UserID,
		SinceVersion: req.SinceVersion,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取通知已读游标版本摘要失败", Data: map[string]any{"userId": req.UserID, "sinceVersion": req.SinceVersion, "err": err.Error()}})
		return nil, err
	}

	cursorVersions := make([]types.NotificationReadCursorVersionItem, 0)
	for _, item := range rpcResp.CursorVersions {
		cursorVersions = append(cursorVersions, types.NotificationReadCursorVersionItem{
			Category: item.Category,
			Version:  item.Version,
		})
	}

	return &types.GetSyncNotificationReadCursorsRes{
		CursorVersions:  cursorVersions,
		MaxVersion:      rpcResp.MaxVersion,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
