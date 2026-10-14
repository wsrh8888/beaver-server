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

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("refresh_token", ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *open_rpc.RefreshTokenReq) (*open_rpc.RefreshTokenRes, error) {
	var token open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("refresh_token_hash = ? AND oauth_id = ? AND revoked_at IS NULL", in.RefreshToken, in.OauthId).First(&token).Error; err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	if time.Now().After(token.RefreshTokenExpiresAt) {
		return nil, errors.New("刷新令牌已过期，请重新授权")
	}

	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	newAccessToken := hex.EncodeToString(accessTokenBytes)

	refreshTokenBytes := make([]byte, 32)
	_, _ = rand.Read(refreshTokenBytes)
	newRefreshToken := hex.EncodeToString(refreshTokenBytes)

	now := time.Now()
	expiresAt := now.Add(2 * time.Hour)
	refreshTokenExpiresAt := now.Add(180 * 24 * time.Hour)

	if err := l.svcCtx.DB.Model(&token).Updates(map[string]interface{}{
		"access_token_hash":        newAccessToken,
		"refresh_token_hash":       newRefreshToken,
		"expires_at":               expiresAt,
		"refresh_token_expires_at": refreshTokenExpiresAt,
	}).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "更新令牌失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, errors.New("刷新令牌失败")
	}

	return &open_rpc.RefreshTokenRes{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    7200,
	}, nil
}
