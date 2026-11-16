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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetEventsByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 按ID拉取通知事件明细
func NewGetEventsByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventsByIdsLogic {
	return &GetEventsByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_events_by_ids", ctx),
	}
}

func (l *GetEventsByIdsLogic) GetEventsByIds(req *types.GetEventsByIdsReq) (resp *types.GetEventsByIdsRes, err error) {
	resp = &types.GetEventsByIdsRes{Events: []types.EventItem{}}

	if len(req.EventIDs) == 0 {
		return resp, nil
	}

	var rows []notification_models.NotificationEvent
	if err = l.svcCtx.DB.WithContext(l.ctx).
		Where("event_id IN ?", req.EventIDs).
		Find(&rows).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询通知事件列表失败",
			Data: map[string]any{"eventIds": req.EventIDs, "err": err.Error()},
		})
		return nil, err
	}

	for _, ev := range rows {
		resp.Events = append(resp.Events, types.EventItem{
			EventID:    ev.EventID,
			EventType:  ev.EventType,
			Category:   ev.Category,
			Version:    ev.Version,
			FromUserID: derefString(ev.FromUserID),
			TargetID:   derefString(ev.TargetID),
			TargetType: ev.TargetType,
			Payload:    string(ev.Payload),
			Priority:   int32(ev.Priority),
			Status:     int32(ev.Status),
			DedupHash:  ev.DedupHash,
			CreatedAt:  time.Time(ev.CreatedAt).UnixMilli(),
			UpdatedAt:  time.Time(ev.UpdatedAt).UnixMilli(),
		})
	}

	return resp, nil
}

func derefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
