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

package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	beaverlog "beaver/utils/beaverlog"
)

type GetOAuthConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetOAuthConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOAuthConfigLogic {
	return &GetOAuthConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("get_oauth_config", ctx),
	}
}

func (l *GetOAuthConfigLogic) GetOAuthConfig(req *types.GetOAuthConfigReq) (resp *types.GetOAuthConfigRes, err error) {
	var app open_models.OpenOAuth
	if err := l.svcCtx.DB.Where("oauth_id = ?", req.OAuthID).First(&app).Error; err != nil {
		return nil, err
	}
	if app.OwnerID != req.UserID {
		return nil, errors.New("无权查看此应用")
	}

	resp = &types.GetOAuthConfigRes{}
	enablePKCE := app.RequirePKCE == 1
	if app.H5 != nil {
		resp.H5Config = &types.H5OAuthConfigInfo{
			Enabled:      app.H5.Enabled,
			RedirectURIs: app.H5.RedirectURIs,
			EnablePKCE:   enablePKCE,
		}
	}
	if app.Desktop != nil && app.Desktop.Enabled {
		authPageURL := ""
		if app.Desktop.CustomScheme != "" {
			authPageURL = fmt.Sprintf("%s/desktop/auth?appId=%s&redirectUri=%s",
				l.svcCtx.Config.OAuth.BaseUrl, req.OAuthID, url.QueryEscape(app.Desktop.CustomScheme))
		}
		resp.DesktopConfig = &types.DesktopOAuthConfigInfo{
			Enabled:      app.Desktop.Enabled,
			CustomScheme: app.Desktop.CustomScheme,
			AuthPageURL:  authPageURL,
		}
	}
	if app.Mobile != nil {
		resp.MobileConfig = &types.MobileOAuthConfigInfo{
			Enabled:            app.Mobile.Enabled,
			IOSBundleID:        app.Mobile.IOSBundleID,
			AndroidPackageName: app.Mobile.AndroidPackageName,
			UniversalLink:      app.Mobile.UniversalLink,
			CustomScheme:       app.Mobile.CustomScheme,
			EnablePKCE:         enablePKCE,
		}
	}
	return resp, nil
}
