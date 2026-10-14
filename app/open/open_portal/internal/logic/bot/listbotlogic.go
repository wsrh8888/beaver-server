/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package bot

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type ListBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBotLogic {
	return &ListBotLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("portal_list_bot", ctx)}
}

func (l *ListBotLogic) ListBot(req *types.ListBotReq) (*types.ListBotRes, error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	query := l.svcCtx.DB.Model(&open_models.OpenBot{}).Where("owner_id = ?", req.UserID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	var rows []open_models.OpenBot
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
		return nil, errors.New("查询失败")
	}

	list := make([]types.BotInfo, 0, len(rows))
	for i := range rows {
		list = append(list, fillBotInfo(l.ctx, l.svcCtx, &rows[i]))
	}
	return &types.ListBotRes{Total: total, List: list}, nil
}
