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

type ResetBotSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewResetBotSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetBotSecretLogic {
	return &ResetBotSecretLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("portal_reset_bot_secret", ctx)}
}

func (l *ResetBotSecretLogic) ResetBotSecret(req *types.ResetBotSecretReq) (*types.ResetBotSecretRes, error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	row, err := loadOwnedBot(l.svcCtx.DB, req.UserID, req.BotID)
	if err != nil {
		return nil, err
	}

	res, err := l.svcCtx.OpenRpc.ResetBotSecret(l.ctx, &open_rpc.ResetBotSecretReq{BotId: row.BotID})
	if err != nil {
		return nil, errors.New("重置密钥失败")
	}
	return &types.ResetBotSecretRes{Secret: res.SignatureSecret}, nil
}
