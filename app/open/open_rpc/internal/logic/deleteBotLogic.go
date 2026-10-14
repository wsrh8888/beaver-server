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
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
)

type DeleteBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBotLogic {
	return &DeleteBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("delete_bot", ctx),
	}
}

func (l *DeleteBotLogic) DeleteBot(in *open_rpc.DeleteBotReq) (*open_rpc.DeleteBotRes, error) {
	if in.BotId == "" {
		return nil, errors.New("bot_id 不能为空")
	}

	result := l.svcCtx.DB.Where("bot_id = ?", in.BotId).Delete(&open_models.OpenBot{})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("机器人不存在")
	}

	_ = l.svcCtx.DB.Where("bot_id = ?", in.BotId).Delete(&open_models.OpenBotMember{})

	return &open_rpc.DeleteBotRes{}, nil
}
