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

type BatchUpdateUserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBatchUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateUserStatusLogic {
	return &BatchUpdateUserStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// BatchUpdateUserStatus 管理后台：批量更新用户状态。
// admin 职责：校验 ids 与 status 合法性，映射为 UpdateUsers 批量改状态 action。
// RPC 职责：UpdateUsers 统一处理状态变更。
func (l *BatchUpdateUserStatusLogic) BatchUpdateUserStatus(req *types.BatchUpdateUserStatusReq) (resp *types.BatchUpdateUserStatusRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("请选择要操作的用户")
	}
	if req.Status < 1 || req.Status > 3 {
		return nil, errors.New("无效的状态值")
	}

	status := int32(req.Status)
	_, err = l.svcCtx.UserRpc.UpdateUsersStatus(l.ctx, &user_rpc.UpdateUsersStatusReq{
		UserIds: req.Ids,
		Status:  status,
	})
	if err != nil {
		l.Errorf("批量更新用户状态失败: %v", err)
		return nil, err
	}
	return &types.BatchUpdateUserStatusRes{}, nil
}
