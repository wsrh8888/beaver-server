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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单
func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuLogic) DeleteMenu(req *types.DeleteMenuReq) (resp *types.DeleteMenuRes, err error) {
	// 检查菜单是否存在
	var menu backend_models.AdminSystemMenu
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			logx.Errorf("菜单不存在: %d", req.Id)
			return nil, err
		}
		logx.Errorf("查询菜单失败: %v", err)
		return nil, err
	}

	// 检查是否有子菜单
	var childCount int64
	err = l.svcCtx.DB.Model(&backend_models.AdminSystemMenu{}).Where("parent_id = ?", req.Id).Count(&childCount).Error
	if err != nil {
		logx.Errorf("检查子菜单失败: %v", err)
		return nil, err
	}

	if childCount > 0 {
		logx.Errorf("无法删除菜单，存在%d个子菜单", childCount)
		return nil, err
	}

	// 删除菜单
	err = l.svcCtx.DB.Delete(&menu).Error
	if err != nil {
		logx.Errorf("删除菜单失败: %v", err)
		return nil, err
	}

	// 删除相关的权限关联
	err = l.svcCtx.DB.Where("menu_id = ?", req.Id).Delete(&backend_models.AdminSystemAuthorityMenu{}).Error
	if err != nil {
		logx.Errorf("删除菜单权限关联失败: %v", err)
		// 这里不返回错误，因为菜单已经删除了
	}

	logx.Infof("菜单删除成功: ID=%d, Name=%s", req.Id, menu.Name)
	return &types.DeleteMenuRes{}, nil
}
