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

type DeleteAppLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除应用
func NewDeleteAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppLogic {
	return &DeleteAppLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAppLogic) DeleteApp(req *types.DeleteAppReq) (resp *types.DeleteAppRes, err error) {

	// 软删除应用（GORM 的 Delete 会设置 deleted_at）
	result := l.svcCtx.DB.Where("app_id = ? AND owner_user_id = ?", req.AppID, req.UserID).Delete(&open_models.OpenApp{})
	if result.Error != nil {
		logx.Errorf("删除应用失败: %v", result.Error)
		return nil, errors.New("删除失败")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("应用不存在或无权限")
	}

	logx.Infof("应用删除成功: app_id=%s", req.AppID)

	return &types.DeleteAppRes{}, nil
}