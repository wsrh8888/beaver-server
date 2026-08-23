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

package event

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEventLogsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEventLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEventLogsLogic {
	return &GetEventLogsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetEventLogsLogic) GetEventLogs(req *types.GetEventLogsReq) (resp *types.GetEventLogsRes, err error) {
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	query := l.svcCtx.DB.Model(&open_models.OpenWebhookLog{}).Where("app_id = ?", req.AppID)
	if req.SubscriptionID > 0 {
		query = query.Where("subscription_id = ?", req.SubscriptionID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询日志失败")
	}

	var logs []open_models.OpenWebhookLog
	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, errors.New("查询日志失败")
	}

	list := make([]types.GetEventLogsResItem, 0, len(logs))
	for _, item := range logs {
		list = append(list, types.GetEventLogsResItem{
			ID:           uint64(item.ID),
			EventID:      item.EventID,
			EventType:    item.EventType,
			TargetURL:    item.TargetURL,
			ResponseCode: item.HTTPStatus,
			CostMs:       item.LatencyMs,
			RetryCount:   item.RetryCount,
			Status:       item.Status,
			ErrorMsg:     item.ErrorMessage,
			CreatedAt:    item.CreatedAt.Unix(),
		})
	}

	return &types.GetEventLogsRes{
		Total: total,
		List:  list,
	}, nil
}
