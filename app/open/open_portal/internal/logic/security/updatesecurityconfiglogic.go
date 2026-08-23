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

type UpdateSecurityConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新安全配置
func NewUpdateSecurityConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSecurityConfigLogic {
	return &UpdateSecurityConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSecurityConfigLogic) UpdateSecurityConfig(req *types.UpdateSecurityConfigReq) (resp *types.UpdateSecurityConfigRes, err error) {

	// 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, err
	}

	// 验证 UserID 是否为应用所有者
	if app.OwnerUserID != req.UserID {
		return nil, errors.New("无权操作此应用")
	}

	// 将 IP 白名单转换为 JSON 字符串
	ipWhitelistJSON := ""
	if len(req.IPWhitelist) > 0 {
		data, _ := json.Marshal(req.IPWhitelist)
		ipWhitelistJSON = string(data)
	}

	// 将 H5 可信域名转换为 JSON 字符串
	trustedDomainsJSON := ""
	if len(req.TrustedDomains) > 0 {
		data, _ := json.Marshal(req.TrustedDomains)
		trustedDomainsJSON = string(data)
	}

	// 更新应用表（只更新 IP 白名单和 H5 可信域名）
	appUpdates := map[string]interface{}{
		"ip_whitelist":    ipWhitelistJSON,
		"trusted_domains": trustedDomainsJSON,
	}

	if err := l.svcCtx.DB.Model(&open_models.OpenApp{}).Where("app_id = ?", req.AppID).Updates(appUpdates).Error; err != nil {
		return nil, err
	}

	return &types.UpdateSecurityConfigRes{}, nil
}
