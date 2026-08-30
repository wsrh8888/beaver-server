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
	"fmt"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)


type CreateAdminUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewCreateAdminUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAdminUserLogic {
	return &CreateAdminUserLogic{logger: logger.New("create_admin_user"), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateAdminUserLogic) CreateAdminUser(req *types.CreateAdminUserReq) (resp *types.CreateAdminUserRes, err error) {
	if req.Phone == "" || req.Password == "" || req.NickName == "" {
		return nil, errors.New("昵称、手机号、密码不能为空")
	}
	var count int64
	if err = l.svcCtx.DB.Model(&backend_models.AdminUser{}).Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("手机号已存在")
	}

	userID := fmt.Sprintf("admin_%d", time.Now().UnixNano())
	adminUser := backend_models.AdminUser{
		UserID:    userID,
		NickName:  req.NickName,
		Phone:     req.Phone,
		Password:  pwd.HahPwd(req.Password),
		Status:    1,
		CreatedBy: req.OperatorID,
	}
	if err = l.svcCtx.DB.Create(&adminUser).Error; err != nil {
		logx.WithContext(l.ctx).Errorf("创建管理员失败: %v", err)
		return nil, err
	}

	if len(req.AuthorityIds) > 0 {
		rows := make([]backend_models.AdminSystemAuthorityUser, 0, len(req.AuthorityIds))
		for _, aid := range req.AuthorityIds {
			rows = append(rows, backend_models.AdminSystemAuthorityUser{
				UserID:      userID,
				AuthorityID: aid,
			})
		}
		if err = l.svcCtx.DB.Create(&rows).Error; err != nil {
			logx.WithContext(l.ctx).Errorf("绑定管理员角色失败: %v", err)
			return nil, err
		}
	}
	l.logger.Info(model.LogMsg{
		Text: "管理员账号创建成功",
		Data: map[string]interface{}{
			"userId":     userID,
			"operatorId": req.OperatorID,
		},
	})
	return &types.CreateAdminUserRes{UserID: userID}, nil
}
