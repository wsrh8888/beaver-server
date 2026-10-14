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

type UpdateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBotLogic {
	return &UpdateBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("update_bot", ctx),
	}
}

func (l *UpdateBotLogic) UpdateBot(in *open_rpc.UpdateBotReq) (*open_rpc.UpdateBotRes, error) {
	if in.BotId == "" {
		return nil, errors.New("botId 不能为空")
	}

	var bot open_models.OpenBot
	if err := l.svcCtx.DB.Where("bot_id = ?", in.BotId).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	updates := map[string]interface{}{}
	if in.Visibility != "" {
		updates["visibility"] = in.Visibility
	}
	if in.SignatureEnabled == 0 || in.SignatureEnabled == 1 {
		updates["signature_enabled"] = in.SignatureEnabled
	}
	if in.Security != nil {
		security := open_models.OpenBotSecurity{
			KeywordsEnabled:    in.Security.KeywordsEnabled,
			IPWhitelistEnabled: in.Security.IpWhitelistEnabled,
		}
		if in.Security.KeywordsEnabled {
			security.Keywords = in.Security.Keywords
		}
		if in.Security.IpWhitelistEnabled {
			security.IPWhitelist = in.Security.IpWhitelist
		}
		updates["security"] = security
	}

	if len(updates) == 0 {
		return &open_rpc.UpdateBotRes{}, nil
	}
	if err := l.svcCtx.DB.Model(&bot).Updates(updates).Error; err != nil {
		return nil, errors.New("更新机器人失败")
	}
	return &open_rpc.UpdateBotRes{}, nil
}
