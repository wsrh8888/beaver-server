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

package auth_public

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GetTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTokenLogic {
	return &GetTokenLogic{
		ctx:    ctx,
		logger: beaverlog.New("get_token", ctx),
		svcCtx: svcCtx,
	}
}

func (l *GetTokenLogic) GetToken(req *types.GetTokenReq) (resp *types.GetTokenRes, err error) {
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", req.AppID, req.AppSecret, 1).First(&app).Error; err != nil {
		l.logger.Warn(model.LogMsg{
			Text: "应用凭证错误",
			Data: map[string]any{"appId": req.AppID},
		})
		return nil, errors.New("应用 ID 或密钥错误")
	}

	var oauthConfig open_models.OpenAppOAuth
	var supportedScopes string
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&oauthConfig).Error; err == nil {
		supportedScopes = oauthConfig.SupportedScopes
	}

	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	accessToken := hex.EncodeToString(accessTokenBytes)

	refreshTokenBytes := make([]byte, 32)
	_, _ = rand.Read(refreshTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	now := time.Now()
	expiresAt := now.Add(2 * time.Hour).Unix()
	refreshTokenExpiresAt := now.Add(180 * 24 * time.Hour).Unix()
	tokenRecord := open_models.OpenOAuthToken{
		AppID:                 req.AppID,
		Token:                 accessToken,
		RefreshToken:          refreshToken,
		ExpiresAt:             expiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		Scope:                 supportedScopes,
		UserID:                "",
	}

	if err := l.svcCtx.DB.Create(&tokenRecord).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "创建访问令牌失败",
			Data: map[string]any{"appId": req.AppID, "err": err.Error()},
		})
		return nil, errors.New("生成令牌失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "客户端凭证令牌签发成功",
		Data: map[string]any{"appId": req.AppID},
	})

	return &types.GetTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    7200,
		TokenType:    "Bearer",
	}, nil
}
