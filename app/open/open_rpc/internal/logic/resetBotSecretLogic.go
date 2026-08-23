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
	uuidUtil "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetBotSecretLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetBotSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetBotSecretLogic {
	return &ResetBotSecretLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetBotSecretLogic) ResetBotSecret(in *open_rpc.ResetBotSecretReq) (*open_rpc.ResetBotSecretRes, error) {
	if in.Id == 0 {
		return nil, errors.New("id 不能为空")
	}

	var bot open_models.OpenBotModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&bot).Error; err != nil {
		return nil, errors.New("机器人不存在")
	}

	security := bot.Security
	security.SignatureEnabled = true
	security.SignatureSecret = uuidUtil.NewV4().String()

	if err := l.svcCtx.DB.Model(&bot).Update("security", security).Error; err != nil {
		return nil, errors.New("重置密钥失败")
	}

	return &open_rpc.ResetBotSecretRes{
		SignatureSecret: security.SignatureSecret,
	}, nil
}
