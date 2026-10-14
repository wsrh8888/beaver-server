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
	"beaver/app/open/open_rpc/types/open_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type DeleteBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewDeleteBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBotLogic {
	return &DeleteBotLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("portal_delete_bot", ctx)}
}

func (l *DeleteBotLogic) DeleteBot(req *types.DeleteBotReq) (*types.DeleteBotRes, error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedBot(l.svcCtx.DB, req.UserID, req.BotID)
	if err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{BotId: row.BotID}); err != nil {
		return nil, errors.New("删除失败")
	}
	return &types.DeleteBotRes{}, nil
}
