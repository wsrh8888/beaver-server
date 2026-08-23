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

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSyncNotificationInboxesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取通知收件箱版本摘要
func NewGetSyncNotificationInboxesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSyncNotificationInboxesLogic {
	return &GetSyncNotificationInboxesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSyncNotificationInboxesLogic) GetSyncNotificationInboxes(req *types.GetSyncNotificationInboxesReq) (resp *types.GetSyncNotificationInboxesRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcResp, err := l.svcCtx.NotificationRpc.GetInboxVersions(l.ctx, &notification_rpc.GetInboxVersionsReq{
		UserId:       req.UserID,
		SinceVersion: req.SinceVersion,
		Limit:        req.Limit,
	})
	if err != nil {
		l.Errorf("获取通知收件箱版本摘要失败: userId=%s, sinceVersion=%d, limit=%d, err=%v", req.UserID, req.SinceVersion, req.Limit, err)
		return nil, err
	}

	inboxVersions := make([]types.NotificationInboxVersionItem, 0)
	for _, item := range rpcResp.InboxVersions {
		inboxVersions = append(inboxVersions, types.NotificationInboxVersionItem{
			EventID: item.EventId,
			Version: item.Version,
		})
	}

	return &types.GetSyncNotificationInboxesRes{
		InboxVersions:   inboxVersions,
		MaxVersion:      rpcResp.MaxVersion,
		ServerTimestamp: time.Now().UnixMilli(),
	}, nil
}
