/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * beaver-server-header-v1
 */

package logic

import (
	"context"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type ListBotsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListBotsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListBotsLogic {
	return &ListBotsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("list_bots", ctx),
	}
}

func (l *ListBotsLogic) ListBots(in *open_rpc.ListBotsReq) (*open_rpc.ListBotsRes, error) {
	page := in.Page
	if page <= 0 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	query := l.svcCtx.DB.Model(&open_models.OpenBot{})
	if in.OwnerUserId != "" {
		query = query.Where("owner_id = ?", in.OwnerUserId)
	}
	if in.Keyword != "" {
		query = query.Where("bot_id LIKE ?", "%"+in.Keyword+"%")
	}
	// status < 0 表示全部（管理端默认 -1）；0 禁用，1 启用
	if in.Status >= 0 {
		query = query.Where("status = ?", in.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var rows []open_models.OpenBot
	offset := int((page - 1) * pageSize)
	if err := query.Order("id DESC").Offset(offset).Limit(int(pageSize)).Find(&rows).Error; err != nil {
		return nil, err
	}

	list := make([]*open_rpc.BotItem, 0, len(rows))
	for i := range rows {
		bot := &rows[i]
		list = append(list, &open_rpc.BotItem{
			Id:               uint32(bot.ID),
			BotId:            bot.BotID,
			OwnerId:          bot.OwnerID,
			Visibility:       bot.Visibility,
			Status:           int32(bot.Status),
			SignatureEnabled: int32(bot.SignatureEnabled),
			CreatedAt:        bot.CreatedAt.UnixMilli(),
		})
	}
	return &open_rpc.ListBotsRes{Total: total, List: list}, nil
}
