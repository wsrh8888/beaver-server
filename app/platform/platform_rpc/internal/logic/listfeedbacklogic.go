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

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListFeedbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFeedbackLogic {
	return &ListFeedbackLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("list_feedback", ctx)}
}

func (l *ListFeedbackLogic) ListFeedback(in *platform_rpc.ListFeedbackReq) (*platform_rpc.ListFeedbackRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&platform_models.FeedbackModel{})
	if in.Status > 0 {
		db = db.Where("status = ?", in.Status)
	}
	if in.Type > 0 {
		db = db.Where("type = ?", in.Type)
	}
	if in.UserId != "" {
		db = db.Where("user_id = ?", in.UserId)
	}
	if in.Keywords != "" {
		db = db.Where("content LIKE ?", "%"+in.Keywords+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "统计反馈失败", Data: map[string]any{"err": err.Error()}})
		return nil, err
	}

	var list []platform_models.FeedbackModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "查询反馈列表失败", Data: map[string]any{"page": page, "pageSize": pageSize, "err": err.Error()}})
		return nil, err
	}

	items := make([]*platform_rpc.FeedbackItem, 0, len(list))
	for _, f := range list {
		items = append(items, toFeedbackItem(f))
	}

	return &platform_rpc.ListFeedbackRes{Total: total, List: items}, nil
}
