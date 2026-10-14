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
	uuidUtil "beaver/utils/uuid"

	beaverlog "beaver/utils/beaverlog"
)

type ResetBotSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewResetBotSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetBotSecretLogic {
	return &ResetBotSecretLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("reset_bot_secret", ctx),
	}
}

func (l *ResetBotSecretLogic) ResetBotSecret(in *open_rpc.ResetBotSecretReq) (*open_rpc.ResetBotSecretRes, error) {
	if in.BotId == "" {
		return nil, errors.New("bot_id 不能为空")
	}

	var bot open_models.OpenBot
	if err := l.svcCtx.DB.Where("bot_id = ?", in.BotId).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	secret := uuidUtil.NewV4().String()
	if err := l.svcCtx.DB.Model(&bot).Updates(map[string]interface{}{
		"signature_enabled": 1,
		"signature_secret": secret,
	}).Error; err != nil {
		return nil, errors.New("重置密钥失败")
	}

	return &open_rpc.ResetBotSecretRes{
		SignatureSecret: secret,
	}, nil
}
