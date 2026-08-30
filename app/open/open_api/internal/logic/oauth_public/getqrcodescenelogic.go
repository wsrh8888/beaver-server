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

package oauth_public

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"beaver/app/open/constants"
	oauthmiddle "beaver/app/open/open_api/internal/middle/oauth"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQrCodeSceneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetQrCodeSceneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQrCodeSceneLogic {
	return &GetQrCodeSceneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetQrCodeSceneLogic) GetQrCodeScene(req *types.GetQrCodeSceneReq) (resp *types.GetQrCodeSceneRes, err error) {
	qrCode, err := l.svcCtx.OAuth.LoadScene(req.SceneID)
	if err != nil {
		return nil, err
	}

	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND status = ?", qrCode.AppID, 1).First(&app).Error; err != nil {
		return nil, err
	}

	expireIn := int64(time.Until(qrCode.ExpiresAt).Seconds())
	if expireIn < 0 {
		expireIn = 0
	}

	var oauthConfig open_models.OpenAppOAuth
	scopeStr := ""
	if err := l.svcCtx.DB.Where("app_id = ?", qrCode.AppID).First(&oauthConfig).Error; err == nil && oauthConfig.SupportedScopes != "" {
		scopeStr = oauthConfig.SupportedScopes
	} else {
		scopes := []string{
			string(constants.ScopeUserProfileRead),
			string(constants.ScopeUserAvatarRead),
		}
		data, _ := json.Marshal(scopes)
		scopeStr = string(data)
	}

	return &types.GetQrCodeSceneRes{
		SceneID:  qrCode.SceneID,
		AppID:    qrCode.AppID,
		AppName:  app.Name,
		AppIcon:  app.Icon,
		Status:   oauthmiddle.QrStatusText(qrCode.Status),
		ExpireIn: expireIn,
		Scopes:   parseScopeList(scopeStr),
	}, nil
}

func parseScopeList(scopeStr string) []string {
	scopeStr = strings.TrimSpace(scopeStr)
	if scopeStr == "" {
		return nil
	}
	if strings.HasPrefix(scopeStr, "[") {
		var scopes []string
		if err := json.Unmarshal([]byte(scopeStr), &scopes); err == nil {
			return scopes
		}
	}
	return strings.FieldsFunc(scopeStr, func(r rune) bool {
		return r == ' ' || r == ','
	})
}
