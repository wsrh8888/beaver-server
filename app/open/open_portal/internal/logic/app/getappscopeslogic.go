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

	"beaver/app/open/constants"
	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppScopesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取应用权限列表
func NewGetAppScopesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppScopesLogic {
	return &GetAppScopesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppScopesLogic) GetAppScopes(req *types.GetAppScopesReq) (resp *types.GetAppScopesRes, err error) {

	// 查询应用
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或无权限")
	}

	// 创建 defaultScopes map 提高查找效率
	defaultMap := make(map[string]bool)
	for _, s := range constants.DefaultScopes {
		defaultMap[string(s)] = true
	}

	scopes := make([]types.ScopeInfo, 0, len(constants.AllScopes))
	for _, scope := range constants.AllScopes {
		scopeStr := string(scope)
		scopes = append(scopes, types.ScopeInfo{
			Scope:       scopeStr,
			Name:        scopeStr,
			Description: constants.ScopeDescription[scope],
			Required:    defaultMap[scopeStr],
		})
	}

	return &types.GetAppScopesRes{
		Scopes: scopes,
	}, nil
}
