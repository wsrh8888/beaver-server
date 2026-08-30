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
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAdminUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdminUserLogic {
	return &UpdateAdminUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UpdateAdminUserLogic) UpdateAdminUser(req *types.UpdateAdminUserReq) (resp *types.UpdateAdminUserRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	var adminUser backend_models.AdminUser
	if err = l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&adminUser).Error; err != nil {
		return nil, errors.New("管理员不存在")
	}

	updates := map[string]interface{}{}
	if req.NickName != "" {
		updates["nick_name"] = req.NickName
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}
	if req.Password != "" {
		updates["password"] = pwd.HahPwd(req.Password)
	}
	if len(updates) > 0 {
		if err = l.svcCtx.DB.Model(&adminUser).Updates(updates).Error; err != nil {
			l.Errorf("更新管理员失败: %v", err)
			return nil, err
		}
	}

	if req.AuthorityIds != nil {
		if err = l.svcCtx.DB.Where("user_id = ?", req.UserID).
			Delete(&backend_models.AdminSystemAuthorityUser{}).Error; err != nil {
			return nil, err
		}
		if len(req.AuthorityIds) > 0 {
			rows := make([]backend_models.AdminSystemAuthorityUser, 0, len(req.AuthorityIds))
			for _, aid := range req.AuthorityIds {
				rows = append(rows, backend_models.AdminSystemAuthorityUser{
					UserID:      req.UserID,
					AuthorityID: aid,
				})
			}
			if err = l.svcCtx.DB.Create(&rows).Error; err != nil {
				l.Errorf("更新管理员角色失败: %v", err)
				return nil, err
			}
		}
	}
	return &types.UpdateAdminUserRes{}, nil
}
