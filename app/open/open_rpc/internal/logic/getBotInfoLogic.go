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

type GetBotInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetBotInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotInfoLogic {
	return &GetBotInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_bot_info", ctx),
	}
}

func (l *GetBotInfoLogic) GetBotInfo(in *open_rpc.GetBotInfoReq) (*open_rpc.GetBotInfoRes, error) {
	var bot open_models.OpenBot
	if err := l.svcCtx.DB.Where("bot_id = ?", in.BotId).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	return &open_rpc.GetBotInfoRes{
		Id:               uint32(bot.ID),
		BotId:            bot.BotID,
		OwnerId:          bot.OwnerID,
		Visibility:       bot.Visibility,
		SignatureEnabled: int32(bot.SignatureEnabled),
		Security: &open_rpc.BotSecurity{
			KeywordsEnabled:    bot.Security.KeywordsEnabled,
			Keywords:           bot.Security.Keywords,
			IpWhitelistEnabled: bot.Security.IPWhitelistEnabled,
			IpWhitelist:        bot.Security.IPWhitelist,
		},
	}, nil
}
