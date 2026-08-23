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
)

type CreateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建菜单
func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMenuLogic) CreateMenu(req *types.CreateMenuReq) (resp *types.CreateMenuRes, err error) {
	// 创建菜单数据
	menu := backend_models.AdminSystemMenu{
		Path:   req.Path,
		Name:   req.Name,
		Hidden: req.Hidden,
		Sort:   req.Sort,
		Title:  req.Title,
		Icon:   req.Icon,
		Status: 1, // 默认启用
	}

	// 处理parent_id
	if req.ParentId != 0 {
		menu.ParentID = &req.ParentId
	}

	// 创建菜单
	err = l.svcCtx.DB.Create(&menu).Error
	if err != nil {
		logx.Errorf("创建菜单失败: %v", err)
		return nil, err
	}

	logx.Infof("菜单创建成功: ID=%d, Name=%s", menu.Id, menu.Name)
	return &types.CreateMenuRes{}, nil
}
