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
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBotInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBotInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotInfoLogic {
	return &GetBotInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetBotInfo 获取推送机器人信息
func (l *GetBotInfoLogic) GetBotInfo(in *open_rpc.GetBotInfoReq) (*open_rpc.GetBotInfoRes, error) {
	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("bot_id = ?", in.BotId).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	return &open_rpc.GetBotInfoRes{
		Id:    uint32(bot.ID),
		BotId: bot.BotID,
		Token: bot.Token,
		Security: &open_rpc.BotSecurity{
			KeywordsEnabled:    bot.Security.KeywordsEnabled,
			Keywords:           bot.Security.Keywords,
			IpWhitelistEnabled: bot.Security.IPWhitelistEnabled,
			IpWhitelist:        bot.Security.IPWhitelist,
			SignatureEnabled:   bot.Security.SignatureEnabled,
			SignatureSecret:    bot.Security.SignatureSecret,
		},
	}, nil
}
