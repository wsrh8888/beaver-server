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
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *types.GetGroupMemberListReq) (resp *types.GetGroupMemberListRes, err error) {
	if req.GroupId == "" {
		return nil, errors.New("群组ID不能为空")
	}

	rpcRes, err := l.svcCtx.GroupRpc.ListGroupMembers(l.ctx, &group_rpc.ListGroupMembersReq{
		GroupId:  req.GroupId,
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		Role:     int32(req.Role),
		Status:   int32(req.Status),
	})
	if err != nil {
		l.Errorf("获取群成员列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		if m.UserId == "" {
			continue
		}
		if _, ok := seen[m.UserId]; ok {
			continue
		}
		seen[m.UserId] = struct{}{}
		userIDs = append(userIDs, m.UserId)
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.GetGroupMemberListItem, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		nick := ""
		if u, ok := users[m.UserId]; ok && u != nil {
			nick = u.NickName
		}
		list = append(list, types.GetGroupMemberListItem{
			Id:              uint(m.Id),
			GroupId:         m.GroupId,
			UserId:          m.UserId,
			MemberNickname:  nick,
			Role:            int(m.Role),
			ProhibitionTime: int(m.ProhibitionMinutes),
			Status:          int(m.Status),
			CreatedAt:       m.CreatedAt,
			UpdatedAt:       m.UpdatedAt,
		})
	}
	return &types.GetGroupMemberListRes{List: list, Total: rpcRes.Total}, nil
}
