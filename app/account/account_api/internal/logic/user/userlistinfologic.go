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

	"beaver/app/account/account_api/internal/svc"
	"beaver/app/account/account_api/internal/types"
	accttypes "beaver/app/account/account_rpc/pb"
)

// UserListInfoLogic 批量查用户信息
//
// 薄转发：HTTP 出入参与同域 RPC 的 proto 类型字段名一致，用 convert 做一次
// JSON 往返即可，不在这里复制业务规则（规则在 account_rpc 里）。
type UserListInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserListInfoLogic 构造 logic。
func NewUserListInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListInfoLogic {
	return &UserListInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserListInfo 批量查用户信息
func (l *UserListInfoLogic) UserListInfo(req *types.UserListInfoReq) (*types.UserListInfoRes, error) {
	rpcReq := &accttypes.UserListInfoReq{UserIdList: req.UserIDList, SinceTimestamp: req.SinceTimestamp}
	rpcRes, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	tmp3 := make(map[string]types.UserLookupInfo, len(rpcRes.UserInfo))
	for k1, v2 := range rpcRes.UserInfo {
		tmp3[k1] = types.UserLookupInfo{UserID: v2.UserId, NickName: v2.NickName, Avatar: v2.Avatar, Version: v2.Version, Email: v2.Email, Abstract: v2.Abstract, Phone: v2.Phone, Status: v2.Status, Source: v2.Source, UserType: v2.UserType, CreatedAt: v2.CreatedAt, UpdatedAt: v2.UpdatedAt}
	}
	return &types.UserListInfoRes{UserInfo: tmp3}, nil
}
