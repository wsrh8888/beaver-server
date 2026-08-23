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

package security

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSecurityConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取安全配置
func NewGetSecurityConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSecurityConfigLogic {
	return &GetSecurityConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSecurityConfigLogic) GetSecurityConfig(req *types.GetSecurityConfigReq) (resp *types.GetSecurityConfigRes, err error) {

	// 查询应用信息
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权查看此应用")
	}

	// 查询安全配置
	var security open_models.OpenAppSecurity
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&security).Error; err != nil {
		// 如果没有配置，返回空配置
		return &types.GetSecurityConfigRes{
			Config: types.SecurityConfigInfo{
				AppID:          app.AppID,
				IPWhitelist:    []string{},
				TrustedDomains: []string{},
			},
		}, nil
	}

	// 解析 IP 白名单
	var ipWhitelist []string
	if security.IPWhitelist != "" {
		json.Unmarshal([]byte(security.IPWhitelist), &ipWhitelist)
	}

	// 解析 H5 可信域名
	var trustedDomains []string
	if security.TrustedDomains != "" {
		json.Unmarshal([]byte(security.TrustedDomains), &trustedDomains)
	}

	return &types.GetSecurityConfigRes{
		Config: types.SecurityConfigInfo{
			AppID:          app.AppID,
			IPWhitelist:    ipWhitelist,
			TrustedDomains: trustedDomains,
		},
	}, nil
}
