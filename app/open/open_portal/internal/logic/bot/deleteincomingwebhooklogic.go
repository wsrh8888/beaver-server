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

package bot

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteIncomingWebhookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteIncomingWebhookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteIncomingWebhookLogic {
	return &DeleteIncomingWebhookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteIncomingWebhookLogic) DeleteIncomingWebhook(req *types.DeleteIncomingWebhookReq) (resp *types.DeleteIncomingWebhookRes, err error) {
	if req.ID == "" {
		return nil, errors.New("id 不能为空")
	}

	botID, err := strconv.ParseUint(req.ID, 10, 64)
	if err != nil || botID == 0 {
		return nil, errors.New("id 无效")
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("id = ?", botID).First(&bot).Error; err != nil {
		return nil, errors.New("记录不存在")
	}
	if bot.AppID == "" {
		return nil, errors.New("无法删除非 Portal 创建的 Bot")
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", bot.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限操作")
	}

	if _, err := l.svcCtx.OpenRpc.DeleteBot(l.ctx, &open_rpc.DeleteBotReq{Id: uint32(bot.ID)}); err != nil {
		return nil, errors.New("删除失败")
	}

	return &types.DeleteIncomingWebhookRes{}, nil
}
