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
	"beaver/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *auth_rpc.ResetPasswordReq) (*auth_rpc.ResetPasswordRes, error) {
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if in.NewPassword == "" {
		return nil, errors.New("新密码不能为空")
	}

	var credential auth_models.AuthCredentialModel
	if err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error; err != nil {
		logx.Errorf("查询用户凭证失败: %v", err)
		return nil, errors.New("用户凭证不存在")
	}

	credential.Password = pwd.HahPwd(in.NewPassword)
	if err := l.svcCtx.DB.Save(&credential).Error; err != nil {
		logx.Errorf("重置密码失败: %v", err)
		return nil, errors.New("重置密码失败")
	}

	logx.Infof("密码重置成功: userID=%s", in.UserId)
	return &auth_rpc.ResetPasswordRes{Success: true}, nil
}
