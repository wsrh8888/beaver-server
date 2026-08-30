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

package user

import (
	"context"
	"errors"

	"beaver/app/open/constants"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.GetUserListReq) (resp *types.GetUserListRes, err error) {
	tokenRecord, err := loadUserAccessToken(l.svcCtx.DB, req.Authorization, constants.ScopeUserProfileRead)
	if err != nil {
		return nil, err
	}
	if len(req.UserIDs) == 0 {
		return nil, errors.New("userIds 不能为空")
	}
	for _, uid := range req.UserIDs {
		if uid != tokenRecord.UserID {
			return nil, errors.New("无权查询该用户信息")
		}
	}

	listRes, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: req.UserIDs,
	})
	if err != nil {
		return nil, err
	}

	scopes := parseScopeList(tokenRecord.Scope)
	userList := make([]types.GetUserListUserItem, 0, len(req.UserIDs))
	for _, uid := range req.UserIDs {
		info, ok := listRes.UserInfo[uid]
		if !ok {
			continue
		}
		item := types.GetUserListUserItem{
			UserID:   info.UserId,
			Nickname: info.NickName,
		}
		if hasScope(scopes, constants.ScopeUserAvatarRead) {
			item.Avatar = info.Avatar
		}
		if hasScope(scopes, constants.ScopeUserEmailRead) {
			item.Email = info.Email
		}
		userList = append(userList, item)
	}

	return &types.GetUserListRes{Users: userList}, nil
}
