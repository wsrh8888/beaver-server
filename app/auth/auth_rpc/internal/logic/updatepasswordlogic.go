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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"beaver/utils/pwd"
)

type UpdatePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewUpdatePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePasswordLogic {
	return &UpdatePasswordLogic{
		ctx:    ctx,
		logger: beaverlog.New("update_password", ctx),
		svcCtx: svcCtx,
	}
}

func (l *UpdatePasswordLogic) UpdatePassword(in *auth_rpc.UpdatePasswordReq) (*auth_rpc.UpdatePasswordRes, error) {
	if in.UserId == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if in.OldPassword == "" || in.NewPassword == "" {
		return nil, errors.New("旧密码和新密码不能为空")
	}

	var credential auth_models.AuthCredentialModel
	err := l.svcCtx.DB.Take(&credential, "user_id = ?", in.UserId).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "查询用户凭证失败",
			Data: map[string]any{"userId": in.UserId, "err": err.Error()},
		})
		return nil, errors.New("用户凭证不存在")
	}

	if !pwd.CheckPad(credential.Password, in.OldPassword) {
		l.logger.Warn(model.LogMsg{
			Text: "旧密码错误",
			Data: map[string]any{"userId": in.UserId},
		})
		return nil, errors.New("旧密码错误")
	}

	credential.Password = pwd.HahPwd(in.NewPassword)
	err = l.svcCtx.DB.Save(&credential).Error
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "更新密码失败",
			Data: map[string]any{"userId": in.UserId, "err": err.Error()},
		})
		return nil, errors.New("更新密码失败")
	}

	l.logger.Info(model.LogMsg{
		Text: "密码修改成功",
		Data: map[string]any{"userId": in.UserId},
	})

	return &auth_rpc.UpdatePasswordRes{Success: true}, nil
}
