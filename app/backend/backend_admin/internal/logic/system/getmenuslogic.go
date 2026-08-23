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

type GetMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取菜单列表
func NewGetMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenusLogic {
	return &GetMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenusLogic) GetMenus(req *types.GetMenuListReq) (resp *types.GetMenuListRes, err error) {
	// 根据用户ID获取用户的权限列表
	var userAuthorities []backend_models.AdminSystemAuthorityUser
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).Find(&userAuthorities).Error
	if err != nil {
		logx.Errorf("查询用户权限失败: %v", err)
		return nil, err
	}

	if len(userAuthorities) == 0 {
		logx.Infof("用户%s没有任何权限", req.UserID)
		return &types.GetMenuListRes{
			List: []types.GetMenuListItem{},
		}, nil
	}

	// 获取所有权限ID
	authorityIDs := make([]uint, len(userAuthorities))
	for i, ua := range userAuthorities {
		authorityIDs[i] = ua.AuthorityID
	}

	// 根据权限ID获取菜单ID列表
	var authorityMenus []backend_models.AdminSystemAuthorityMenu
	err = l.svcCtx.DB.Where("authority_id IN ?", authorityIDs).Find(&authorityMenus).Error
	if err != nil {
		logx.Errorf("查询权限菜单关联失败: %v", err)
		return nil, err
	}

	if len(authorityMenus) == 0 {
		logx.Infof("用户%s的权限没有任何菜单", req.UserID)
		return &types.GetMenuListRes{
			List: []types.GetMenuListItem{},
		}, nil
	}

	// 获取所有菜单ID（去重）
	menuIDMap := make(map[uint]bool)
	menuIDs := make([]uint, 0)
	for _, am := range authorityMenus {
		if !menuIDMap[am.MenuID] {
			menuIDMap[am.MenuID] = true
			menuIDs = append(menuIDs, am.MenuID)
		}
	}

	// 查询菜单详情
	var menus []backend_models.AdminSystemMenu
	err = l.svcCtx.DB.Where("id IN ? AND status = ?", menuIDs, 1).Order("sort asc").Find(&menus).Error
	if err != nil {
		logx.Errorf("查询菜单详情失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.GetMenuListItem
	for _, menu := range menus {
		parentId := uint(0)
		if menu.ParentID != nil {
			parentId = *menu.ParentID
		}

		list = append(list, types.GetMenuListItem{
			Id:       menu.Id,
			ParentId: parentId,
			Path:     menu.Path,
			Name:     menu.Name,
			Hidden:   menu.Hidden,
			Sort:     menu.Sort,
			Title:    menu.Title,
			Icon:     menu.Icon,
		})
	}

	return &types.GetMenuListRes{
		List: list,
	}, nil
}
