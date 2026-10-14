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
	"fmt"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"
	uuidUtil "beaver/utils/uuid"

	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

const botPublicWebhookPath = "/api/open/bot_public/v1/send"

type CreateBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateBotLogic {
	return &CreateBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("create_bot", ctx),
	}
}

func (l *CreateBotLogic) CreateBot(in *open_rpc.CreateBotReq) (*open_rpc.CreateBotRes, error) {
	if in.BotId == "" {
		return nil, fmt.Errorf("bot_id 不能为空")
	}
	if in.OwnerId == "" {
		return nil, fmt.Errorf("owner_id 不能为空")
	}

	visibility := in.Visibility
	if visibility == "" {
		visibility = open_models.BotVisibilityPrivate
	}

	webhookToken := uuidUtil.NewV4().String()
	signatureSecret := uuidUtil.NewV4().String()

	bot := &open_models.OpenBot{
		OwnerID:          in.OwnerId,
		BotID:            in.BotId,
		Token:            webhookToken,
		Visibility:       visibility,
		SignatureEnabled: 1,
		SignatureSecret:  signatureSecret,
		Security:         open_models.OpenBotSecurity{},
	}

	if err := l.svcCtx.DB.Create(bot).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "创建 Bot 记录失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, fmt.Errorf("创建 Bot 记录失败")
	}

	return &open_rpc.CreateBotRes{
		Id:              uint32(bot.ID),
		BotId:           in.BotId,
		WebhookUrl:      botPublicWebhookPath,
		Token:           webhookToken,
		SignatureSecret: signatureSecret,
	}, nil
}
