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

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
)


type RevokeTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewRevokeTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeTokenLogic {
	return &RevokeTokenLogic{
		ctx:    ctx,
		logger: logger.New("revoke_token"),
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

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND app_secret = ? AND status = ?", appID, appSecret, 1).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或凭证错误")
	}

	var tokenRecord open_models.OpenOAuthToken
	if err := l.svcCtx.DB.Where("token = ? OR refresh_token = ?", req.Token, req.Token).First(&tokenRecord).Error; err != nil {
		return nil, errors.New("令牌不存在")
	}
	if tokenRecord.AppID != appID {
		return nil, errors.New("无权撤销该令牌")
	}

	result := l.svcCtx.DB.Where("id = ?", tokenRecord.ID).Delete(&open_models.OpenOAuthToken{})
	if result.Error != nil {
		return nil, errors.New("撤销令牌失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "OAuth撤销Token成功",
		Data: map[string]interface{}{
			"appId": appID,
		},
	})

	return &types.RevokeTokenRes{Success: true}, nil
}
