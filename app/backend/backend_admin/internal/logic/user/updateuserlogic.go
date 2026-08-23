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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// UpdateUser 管理后台：更新用户。
// admin 职责：校验 userId，将可选字段组装为 patch 语义。
// RPC 职责：UpdateUsers 处理资料更新，UpdateUsersStatus 处理状态变更。
func (l *UpdateUserLogic) UpdateUser(req *types.UpdateUserReq) (resp *types.UpdateUserRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	rpcReq := &user_rpc.UpdateUsersReq{
		UserIds: []string{req.UserID},
	}
	if req.NickName != nil {
		rpcReq.PatchNickName = req.NickName
	}
	if req.Email != nil {
		rpcReq.PatchEmail = req.Email
	}
	if req.Avatar != nil {
		rpcReq.PatchAvatar = req.Avatar
	}
	if req.Abstract != nil {
		rpcReq.PatchAbstract = req.Abstract
	}

	hasPatch := req.NickName != nil || req.Email != nil || req.Avatar != nil || req.Abstract != nil
	if hasPatch {
		_, err = l.svcCtx.UserRpc.UpdateUsers(l.ctx, rpcReq)
		if err != nil {
			l.Errorf("更新用户失败: %v", err)
			return nil, err
		}
	}

	if req.Status != nil {
		_, err = l.svcCtx.UserRpc.UpdateUsersStatus(l.ctx, &user_rpc.UpdateUsersStatusReq{
			UserIds: []string{req.UserID},
			Status:  int32(*req.Status),
		})
		if err != nil {
			l.Errorf("更新用户状态失败: %v", err)
			return nil, err
		}
	}
	return &types.UpdateUserRes{}, nil
}
