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

package app

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用详情
func NewGetAppDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppDetailLogic {
	return &GetAppDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppDetailLogic) GetAppDetail(req *types.GetAppDetailReq) (resp *types.GetAppDetailRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	// 查询应用详情
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限访问")
	}

	// 2. 对 AppSecret 进行掩码处理（只显示前8位和后8位）
	maskedSecret := maskSecret(app.AppSecret)

	return &types.GetAppDetailRes{
		App: types.AppInfo{
			AppID:       app.AppID,
			Name:        app.Name,
			Description: app.Description,
			Icon:        app.Icon,
			AppSecret:   maskedSecret,
			Status:      app.Status,
			// 能力开关
			EnableRobot:   app.EnableRobot,
			EnableOAuth:   app.EnableOAuth,
			EnableWebhook: app.EnableWebhook,
			CreatedAt:     app.CreatedAt.Unix(),
		},
	}, nil
}

// maskSecret 对密钥进行掩码处理
func maskSecret(secret string) string {
	if len(secret) <= 16 {
		return "****"
	}
	return secret[:8] + "****" + secret[len(secret)-8:]
}
