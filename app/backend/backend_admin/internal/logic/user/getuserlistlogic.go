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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetUserList 管理后台：用户列表查询。
// admin 职责：运营筛选条件映射、响应字段适配前端协议。
// RPC 职责：ListUsers 领域查询，不与本 HTTP 接口 1:1。
func (l *GetUserListLogic) GetUserList(req *types.GetUserListReq) (resp *types.GetUserListRes, err error) {
	rpcRes, err := l.svcCtx.UserRpc.ListUsers(l.ctx, &user_rpc.ListUsersReq{
		Page:     int32(req.Page),
		PageSize: int32(req.PageSize),
		Email:    req.Email,
		Keyword:  req.Keyword,
		Status:   int32(req.Status),
		Source:   int32(req.Source),
		UserType: int32(req.UserType),
	})
	if err != nil {
		l.Errorf("获取用户列表失败: %v", err)
		return nil, err
	}

	list := make([]types.UserInfo, 0, len(rpcRes.List))
	for _, u := range rpcRes.List {
		list = append(list, types.UserInfo{
			Id: u.UserId, NickName: u.NickName, Email: u.Email, Abstract: u.Abstract,
			Avatar: u.Avatar, Status: int(u.Status), Source: int(u.Source), UserType: int(u.UserType),
			CreateTime: u.CreatedAt, UpdateTime: u.UpdatedAt,
		})
	}
	return &types.GetUserListRes{List: list, Total: rpcRes.Total}, nil
}
