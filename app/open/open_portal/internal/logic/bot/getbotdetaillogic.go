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

	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type GetBotDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetBotDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotDetailLogic {
	return &GetBotDetailLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("portal_bot_detail", ctx)}
}

func (l *GetBotDetailLogic) GetBotDetail(req *types.GetBotDetailReq) (*types.GetBotDetailRes, error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedBot(l.svcCtx.DB, req.UserID, req.BotID)
	if err != nil {
		return nil, err
	}
	info := fillBotInfo(l.ctx, l.svcCtx, row)
	return &types.GetBotDetailRes{Bot: info}, nil
}
