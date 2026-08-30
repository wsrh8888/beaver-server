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
	"fmt"
	"time"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBotDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取机器人详情
func NewGetBotDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotDetailLogic {
	return &GetBotDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBotDetailLogic) GetBotDetail(req *types.GetBotDetailReq) (resp *types.GetBotDetailRes, err error) {
	// 1. 查询群内机器人信息（通过 bot_id 查询）
	var bot group_models.GroupBotModel
	if err := l.svcCtx.DB.Where("bot_id = ?", req.BotID).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	// 2. 通过 user_rpc 获取用户基础信息（昵称、头像）
	userRes, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user_rpc.UserInfoReq{
		UserID: bot.BotID,
	})
	if err != nil || userRes.UserInfo == nil {
		return nil, err
	}

	// 3. 通过 open_rpc 获取 Webhook Token 和安全设置
	botInfoRes, err := l.svcCtx.OpenRpc.GetBotInfo(l.ctx, &open_rpc.GetBotInfoReq{
		BotId: bot.BotID,
	})
	if err != nil {
		return nil, err
	}

	// 4. 拼接完整 Webhook URL
	fullWebhookURL := fmt.Sprintf("%s/api/open/bot_public/v1/send?token=%s", l.svcCtx.Config.Domain, botInfoRes.Token)

	return &types.GetBotDetailRes{
		BotID:         bot.BotID,
		Name:          userRes.UserInfo.NickName,
		Description:   userRes.UserInfo.Abstract, // 用户简介
		Avatar:        userRes.UserInfo.Avatar,
		WebhookURL:    fullWebhookURL,
		Type:          bot.Type,
		Status:        bot.Status,
		CreatorUserID: bot.CreatorID,
		CreatedAt:     time.Time(bot.CreatedAt).Unix(),
		Security: types.BotSecurity{
			KeywordsEnabled:    botInfoRes.Security.KeywordsEnabled,
			Keywords:           botInfoRes.Security.Keywords,
			IPWhitelistEnabled: botInfoRes.Security.IpWhitelistEnabled,
			IPWhitelist:        botInfoRes.Security.IpWhitelist,
			SignatureEnabled:   botInfoRes.Security.SignatureEnabled,
			SignatureSecret:    botInfoRes.Security.SignatureSecret,
		},
	}, nil
}
