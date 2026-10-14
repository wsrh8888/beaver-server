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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ExchangeTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewExchangeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExchangeTokenLogic {
	return &ExchangeTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("exchange_token", ctx),
	}
}

func (l *ExchangeTokenLogic) ExchangeToken(in *open_rpc.ExchangeTokenReq) (*open_rpc.ExchangeTokenRes, error) {
	var app open_models.OpenOAuth
	if err := l.svcCtx.DB.Where("oauth_id = ? AND status = ?", in.OauthId, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	var authCode open_models.OpenOAuthCode
	if err := l.svcCtx.DB.Where("code_hash = ? AND oauth_id = ? AND used_at IS NULL", in.Code, in.OauthId).First(&authCode).Error; err != nil {
		return nil, errors.New("授权码无效或已使用")
	}

	if time.Now().After(authCode.ExpiresAt) {
		return nil, errors.New("授权码已过期")
	}

	if err := l.svcCtx.DB.Model(&authCode).Update("used_at", time.Now()).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "更新授权码状态失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, errors.New("处理授权码失败")
	}

	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	accessToken := hex.EncodeToString(accessTokenBytes)

	refreshTokenBytes := make([]byte, 32)
	_, _ = rand.Read(refreshTokenBytes)
	refreshToken := hex.EncodeToString(refreshTokenBytes)

	now := time.Now()
	accessTTL := app.AccessTokenTTL
	if accessTTL <= 0 {
		accessTTL = 7200
	}
	refreshTTL := app.RefreshTokenTTL
	if refreshTTL <= 0 {
		refreshTTL = 2592000
	}
	expiresAt := now.Add(time.Duration(accessTTL) * time.Second)
	refreshTokenExpiresAt := now.Add(time.Duration(refreshTTL) * time.Second)
	tokenRecord := open_models.OpenOAuthToken{
		OAuthID:               in.OauthId,
		Token:                 accessToken,
		RefreshToken:          refreshToken,
		ExpiresAt:             expiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
		Scope:                 authCode.Scope,
		UserID:                authCode.UserID,
	}

	if err := l.svcCtx.DB.Create(&tokenRecord).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "创建访问令牌失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, errors.New("生成令牌失败")
	}

	return &open_rpc.ExchangeTokenRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTTL),
		TokenType:    "Bearer",
		UserId:       authCode.UserID,
	}, nil
}
