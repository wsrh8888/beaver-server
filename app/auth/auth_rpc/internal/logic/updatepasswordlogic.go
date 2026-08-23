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

package logic

import (
	"context"
	"errors"

	"beaver/app/auth/auth_models"
	"beaver/app/auth/auth_rpc/internal/svc"
	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/utils/logger"
	"beaver/utils/logger/model"
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)


type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		logger: logger.New("update_password"),
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(in *auth_rpc.UpdatePasswordReq) (*auth_rpc.UpdatePasswordRes, error) {
	// 验证必填字段
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if in.OldPassword == "" || in.NewPassword == "" {
		return nil, errors.New("旧密码和新密码不能为空")
	}

	// 查询用户凭证
	var credential auth_models.AuthCredentialModel
	err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error
	if err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return nil, errors.New("用户凭证不存在")
	}

	// 验证旧密码
	if !pwd.CheckPad(credential.Password, in.OldPassword) {
		return nil, errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword := pwd.HahPwd(in.NewPassword)

	// 更新密码
	credential.Password = hashedPassword
	err = l.svcCtx.DB.Save(&credential).Error
	if err != nil {
		logx.Errorf("更新密码失败: %v", err)
		return nil, errors.New("更新密码失败")
	}

	logx.Infof("密码更新成功: userID=%s", in.UserId)
	l.logger.Info(model.LogMsg{
		Text: "密码修改成功",
		Data: map[string]interface{}{
			"userId": in.UserId,
		},
	})

	return &auth_rpc.UpdatePasswordRes{
		Success: true,
	}, nil
}
