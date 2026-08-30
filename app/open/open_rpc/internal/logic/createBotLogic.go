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
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	uuidUtil "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBotLogic {
	return &CreateBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateBot 创建推送机器人（open 服务负责生成 Token、安全设置等）
func (l *CreateBotLogic) CreateBot(in *open_rpc.CreateBotReq) (*open_rpc.CreateBotRes, error) {
	// 1. 使用传入的机器人用户 ID（由 user_rpc 生成）
	if in.BotId == "" {
		return nil, fmt.Errorf("bot_id 不能为空")
	}

	// 2. 生成 Webhook Token（用于身份验证）
	webhookToken := uuidUtil.NewV4().String()

	// 3. 生成签名密钥（默认启用签名校验）
	signatureSecret := uuidUtil.NewV4().String()

	// 4. 生成 Webhook URL
	webhookURL := fmt.Sprintf("/api/webhook/%s", webhookToken)

	// 5. 创建 open_bots 记录（默认启用签名校验）
	bot := &open_models.OpenBotModel{
		BotID:   in.BotId, // 使用传入的用户 ID
		GroupID: in.GroupId,
		Token:   webhookToken,
		Status:  1,
		Security: open_models.OpenBotSecurity{
			SignatureEnabled:   true, // 默认启用签名校验
			SignatureSecret:    signatureSecret,
			KeywordsEnabled:    false, // 关键词校验默认关闭
			IPWhitelistEnabled: false, // IP白名单默认关闭
		},
	}

	if err := l.svcCtx.DB.Create(bot).Error; err != nil {
		logx.Errorf("创建 Bot 记录失败: %v", err)
		return nil, fmt.Errorf("创建 Bot 记录失败")
	}

	return &open_rpc.CreateBotRes{
		Id:         uint32(bot.ID),
		BotId:      in.BotId,
		WebhookUrl: webhookURL,
		Token:      webhookToken,
	}, nil
}
