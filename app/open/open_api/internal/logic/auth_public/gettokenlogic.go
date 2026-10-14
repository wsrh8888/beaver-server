/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
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
	if req.AppKey == "" || req.AppSecret == "" {
		return nil, errors.New("appKey 或 appSecret 错误")
	}

	var partner open_models.OpenDigital
	if err := l.svcCtx.DB.Where(
		"app_key = ? AND app_secret_hash = ? AND runtime = ?",
		req.AppKey, req.AppSecret, open_models.DigitalRuntimeExternal,
	).First(&partner).Error; err != nil {
		l.logger.Warn(model.LogMsg{
			Text: "数字伙伴凭证错误",
			Data: map[string]any{"appKey": req.AppKey},
		})
		return nil, errors.New("appKey 或 appSecret 错误")
	}

	if partner.Status != open_models.DigitalStatusEnabled {
		return nil, errors.New("数字伙伴已禁用")
	}

	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	accessToken := hex.EncodeToString(accessTokenBytes)

	now := time.Now()
	expiresAt := now.Add(2 * time.Hour)
	tokenRecord := open_models.OpenDigitalToken{
		DigitalID: partner.DigitalID,
		Token:     accessToken,
		ExpiresAt: expiresAt,
	}

	if err := l.svcCtx.DB.Create(&tokenRecord).Error; err != nil {
		l.logger.Error(model.LogMsg{
			Text: "创建访问令牌失败",
			Data: map[string]any{"appKey": req.AppKey, "err": err.Error()},
		})
		return nil, errors.New("生成令牌失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "外接数字伙伴令牌签发成功",
		Data: map[string]any{"appKey": req.AppKey},
	})

	return &types.GetTokenRes{
		AccessToken: accessToken,
		ExpiresIn:   7200,
		TokenType:   "Bearer",
	}, nil
}
