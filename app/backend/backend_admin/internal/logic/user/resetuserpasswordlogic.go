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

	"beaver/app/auth/auth_rpc/types/auth_rpc"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type ResetUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{logger: logger.New("reset_user_password"), ctx: ctx, svcCtx: svcCtx}
}

// ResetUserPassword 管理后台：重置用户密码。
// admin 职责：校验 userId/newPassword，调用 AuthRpc 重置凭证（密码属认证域，不走 UserRpc）。
// RPC 职责：AuthRpc.ResetPassword 管理密码哈希与凭证状态。
func (l *ResetUserPasswordLogic) ResetUserPassword(req *types.ResetUserPasswordReq) (resp *types.ResetUserPasswordRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}
	if req.NewPassword == "" {
		return nil, errors.New("新密码不能为空")
	}

	_, err = l.svcCtx.AuthRpc.ResetPassword(l.ctx, &auth_rpc.ResetPasswordReq{
		UserId:      req.UserID,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		logx.WithContext(l.ctx).Errorf("重置用户密码失败: %v", err)
		return nil, err
	}
	l.logger.Info(model.LogMsg{
		Text: "管理员重置用户密码成功",
		Data: map[string]interface{}{
			"userId": req.UserID,
		},
	})
	return &types.ResetUserPasswordRes{}, nil
}
