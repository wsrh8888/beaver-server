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

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type ResetAppSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置应用密钥
func NewResetAppSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetAppSecretLogic {
	return &ResetAppSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAppSecretLogic) ResetAppSecret(req *types.ResetAppSecretReq) (resp *types.ResetAppSecretRes, err error) {

	// 生成新密钥
	newSecret := uuid.New().String() + uuid.New().String()

	// 更新密钥
	result := l.svcCtx.DB.Model(&open_models.OpenApp{}).
		Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).
		Update("app_secret", newSecret)

	if result.Error != nil {
		logx.Errorf("重置密钥失败: %v", result.Error)
		return nil, errors.New("重置失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("应用不存在或无权限")
	}

	logx.Infof("应用密钥重置成功: app_id=%s", req.AppID)

	return &types.ResetAppSecretRes{
		AppSecret: newSecret,
	}, nil
}
