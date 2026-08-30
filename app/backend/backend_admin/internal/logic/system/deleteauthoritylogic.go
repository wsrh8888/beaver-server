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

package system

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAuthorityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAuthorityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAuthorityLogic {
	return &DeleteAuthorityLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteAuthorityLogic) DeleteAuthority(req *types.DeleteAuthorityReq) (resp *types.DeleteAuthorityRes, err error) {
	if req.Id == 0 {
		return nil, errors.New("角色ID不能为空")
	}
	var userCount int64
	if err = l.svcCtx.DB.Model(&backend_models.AdminSystemAuthorityUser{}).
		Where("authority_id = ?", req.Id).Count(&userCount).Error; err != nil {
		return nil, err
	}
	if userCount > 0 {
		return nil, errors.New("该角色仍有关联管理员，无法删除")
	}
	if err = l.svcCtx.DB.Where("authority_id = ?", req.Id).
		Delete(&backend_models.AdminSystemAuthorityMenu{}).Error; err != nil {
		return nil, err
	}
	if err = l.svcCtx.DB.Delete(&backend_models.AdminSystemAuthority{}, req.Id).Error; err != nil {
		l.Errorf("删除角色失败: %v", err)
		return nil, err
	}
	return &types.DeleteAuthorityRes{}, nil
}
