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

	"gorm.io/gorm"
)

type DeleteMenuLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单
func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		logger: beaverlog.New("delete_menu", ctx),
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
			l.logger.Error(model.LogMsg{
				Text: "菜单不存在",
				Data: map[string]interface{}{"menuId": req.Id},
			})
			return nil, err
		}
		l.logger.Error(model.LogMsg{
			Text: "查询菜单失败",
			Data: map[string]interface{}{"menuId": req.Id, "err": err.Error()},
		})
		return nil, err
	}

	// 检查是否有子菜单
	var childCount int64
	err = l.svcCtx.DB.Model(&backend_models.AdminSystemMenu{}).Where("parent_id = ?", req.Id).Count(&childCount).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "检查子菜单失败",
			Data: map[string]interface{}{"menuId": req.Id, "err": err.Error()},
		})
		return nil, err
	}

	if childCount > 0 {
		l.logger.Error(model.LogMsg{
			Text: "无法删除菜单：存在子菜单",
			Data: map[string]interface{}{"menuId": req.Id, "childCount": childCount},
		})
		return nil, err
	}

	// 删除菜单
	err = l.svcCtx.DB.Delete(&menu).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "删除菜单失败",
			Data: map[string]interface{}{"menuId": req.Id, "err": err.Error()},
		})
		return nil, err
	}

	// 删除相关的权限关联
	err = l.svcCtx.DB.Where("menu_id = ?", req.Id).Delete(&backend_models.AdminSystemAuthorityMenu{}).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "删除菜单权限关联失败",
			Data: map[string]interface{}{"menuId": req.Id, "err": err.Error()},
		})
		// 这里不返回错误，因为菜单已经删除了
	}

	l.logger.Info(model.LogMsg{
		Text: "菜单删除成功",
		Data: map[string]interface{}{"menuId": req.Id, "menuName": menu.Name},
	})
	return &types.DeleteMenuRes{}, nil
}
