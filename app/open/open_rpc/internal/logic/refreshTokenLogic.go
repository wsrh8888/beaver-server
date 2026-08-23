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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshTokenLogic) RefreshToken(in *open_rpc.RefreshTokenReq) (*open_rpc.RefreshTokenRes, error) {
	// 1. 查询 Refresh Token
	var token open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("refresh_token = ? AND app_id = ?", in.RefreshToken, in.AppId).First(&token).Error; err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 2. 检查 Refresh Token 是否过期
	if time.Now().Unix() > token.RefreshTokenExpiresAt {
		return nil, errors.New("刷新令牌已过期，请重新授权")
	}

	// 3. 生成新的 Access Token
	accessTokenBytes := make([]byte, 32)
	_, _ = rand.Read(accessTokenBytes)
	newAccessToken := hex.EncodeToString(accessTokenBytes)

	// 4. 生成新的 Refresh Token
	refreshTokenBytes := make([]byte, 32)
	_, _ = rand.Read(refreshTokenBytes)
	newRefreshToken := hex.EncodeToString(refreshTokenBytes)

	// 5. 更新数据库
	now := time.Now()
	expiresAt := now.Add(2 * time.Hour).Unix()                    // access_token 2小时过期
	refreshTokenExpiresAt := now.Add(180 * 24 * time.Hour).Unix() // refresh_token 180天过期

	if err := l.svcCtx.DB.Model(&token).Updates(map[string]interface{}{
		"token":                    newAccessToken,
		"refresh_token":            newRefreshToken,
		"expires_at":               expiresAt,
		"refresh_token_expires_at": refreshTokenExpiresAt,
	}).Error; err != nil {
		logx.Errorf("更新令牌失败: %v", err)
		return nil, errors.New("刷新令牌失败")
	}

	return &open_rpc.RefreshTokenRes{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    7200, // 2小时
	}, nil
}
