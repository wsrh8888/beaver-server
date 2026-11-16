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

package user_machine

import (
	"context"

	"beaver/app/account/account_api/internal/svc"
	"beaver/app/account/account_api/internal/types"
	accttypes "beaver/app/account/account_rpc/pb"
	"beaver/common/convert"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserUpdateStatusLogic 批量修改用户状态
//
// 薄转发：HTTP 出入参与同域 RPC 的 proto 类型字段名一致，用 convert 做一次
// JSON 往返即可，不在这里复制业务规则（规则在 account_rpc 里）。
type UserUpdateStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserUpdateStatusLogic 构造 logic。
func NewUserUpdateStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateStatusLogic {
	return &UserUpdateStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserUpdateStatus 批量修改用户状态
func (l *UserUpdateStatusLogic) UserUpdateStatus(req *types.UserUpdateStatusReq) (*types.UserUpdateStatusRes, error) {
	rpcReq, err := convert.To[accttypes.UpdateUsersStatusReq](req)
	if err != nil {
		return nil, err
	}
	rpcRes, err := l.svcCtx.UserRpc.UpdateUsersStatus(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return convert.To[types.UserUpdateStatusRes](rpcRes)
}
