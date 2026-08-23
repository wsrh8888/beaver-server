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

type UpdateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBotLogic {
	return &UpdateBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateBotLogic) UpdateBot(in *open_rpc.UpdateBotReq) (*open_rpc.UpdateBotRes, error) {
	if in.BotId == "" {
		return nil, errors.New("botId 不能为空")
	}
	if in.Security == nil {
		return &open_rpc.UpdateBotRes{}, nil
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("bot_id = ? AND status = 1", in.BotId).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	security := open_models.OpenBotSecurity{
		KeywordsEnabled:    in.Security.KeywordsEnabled,
		IPWhitelistEnabled: in.Security.IpWhitelistEnabled,
		SignatureEnabled:   in.Security.SignatureEnabled,
	}
	if in.Security.KeywordsEnabled {
		security.Keywords = in.Security.Keywords
	}
	if in.Security.IpWhitelistEnabled {
		security.IPWhitelist = in.Security.IpWhitelist
	}
	if in.Security.SignatureSecret != "" {
		security.SignatureSecret = in.Security.SignatureSecret
	} else {
		security.SignatureSecret = bot.Security.SignatureSecret
	}

	if err := l.svcCtx.DB.Model(&bot).Update("security", security).Error; err != nil {
		return nil, errors.New("更新机器人安全设置失败")
	}

	return &open_rpc.UpdateBotRes{}, nil
}
