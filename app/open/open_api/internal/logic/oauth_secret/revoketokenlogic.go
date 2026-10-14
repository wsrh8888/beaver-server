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

package oauth_secret

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type RevokeTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewRevokeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeTokenLogic {
	return &RevokeTokenLogic{
		ctx:    ctx,
		logger: beaverlog.New("revoke_token", ctx),
		svcCtx: svcCtx,
	}
}

func (l *RevokeTokenLogic) RevokeToken(req *types.RevokeTokenReq, appID, appSecret string) (resp *types.RevokeTokenRes, err error) {
	if req.Token == "" {
		return nil, errors.New("token 不能为空")
	}
	if appID == "" || appSecret == "" {
		return nil, errors.New("应用凭证不能为空")
	}

	var app open_models.OpenOAuth
	if err := l.svcCtx.DB.Where("oauth_id = ? AND oauth_secret_hash = ? AND status = ?", appID, appSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或凭证错误")
	}

	var tokenRecord open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("access_token_hash = ? OR refresh_token_hash = ?", req.Token, req.Token).First(&tokenRecord).Error; err != nil {
		return nil, errors.New("令牌不存在")
	}
	if tokenRecord.OAuthID != appID {
		return nil, errors.New("无权撤销该令牌")
	}

	now := time.Now()
	if err := l.svcCtx.DB.Model(&tokenRecord).Update("revoked_at", now).Error; err != nil {
		return nil, errors.New("撤销令牌失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "OAuth撤销Token成功",
		Data: map[string]interface{}{
			"oauthId": appID,
		},
	})

	return &types.RevokeTokenRes{Success: true}, nil
}
