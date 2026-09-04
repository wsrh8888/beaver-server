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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CreateMenuLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建菜单
func NewCreateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMenuLogic {
	return &CreateMenuLogic{
		logger: beaverlog.New("create_menu", ctx),
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
		l.logger.Error(model.LogMsg{
			Text: "创建菜单失败",
			Data: map[string]interface{}{"err": err.Error()},
		})
		return nil, err
	}

	l.logger.Info(model.LogMsg{
		Text: "菜单创建成功",
		Data: map[string]interface{}{"menuId": menu.Id, "menuName": menu.Name},
	})
	return &types.CreateMenuRes{}, nil
}
