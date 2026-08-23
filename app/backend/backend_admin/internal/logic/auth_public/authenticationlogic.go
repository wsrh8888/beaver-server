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

package auth_public

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/utils/jwts"
	utils "beaver/utils/list"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 管理员 Token 认证
func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationReq) (resp *types.AuthenticationRes, err error) {
	if utils.InListByRegex(l.svcCtx.Config.WhiteList, req.ValidPath) {
		return &types.AuthenticationRes{}, nil
	}

	claims, err := jwts.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("认证失败")
	}

	var adminUser backend_models.AdminUser
	err = l.svcCtx.DB.Take(&adminUser, "user_id = ? AND status = ?", claims.UserID, 1).Error
	if err != nil {
		return nil, errors.New("管理员用户不存在或已被禁用")
	}

	key := fmt.Sprintf("admin_login_%s", claims.UserID)
	token, _ := l.svcCtx.Redis.Get(key).Result()
	if token != req.Token {
		return nil, errors.New("token已失效")
	}

	return &types.AuthenticationRes{
		UserID: claims.UserID,
	}, nil
}
