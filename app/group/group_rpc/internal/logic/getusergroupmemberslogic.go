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

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserGroupMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserGroupMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserGroupMembersLogic {
	return &GetUserGroupMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserGroupMembersLogic) GetUserGroupMembers(in *group_rpc.GetUserGroupMembersReq) (*group_rpc.GetUserGroupMembersRes, error) {
	// 获取用户加入的所有群组
	var userMembers []group_models.GroupMemberModel

	err := l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ? AND status = ?", in.UserID, 1).
		Find(&userMembers).Error

	if err != nil {
		l.Errorf("查询用户群组失败: %v", err)
		return nil, err
	}

	// 提取用户加入的群组ID
	groupIDs := make([]string, 0, len(userMembers))
	for _, member := range userMembers {
		groupIDs = append(groupIDs, member.GroupID)
	}

	if len(groupIDs) == 0 {
		l.Infof("用户未加入任何群组，用户ID: %s", in.UserID)
		return &group_rpc.GetUserGroupMembersRes{
			MemberIDs: []string{},
		}, nil
	}

	// 获取所有群组的成员列表
	var allMembers []group_models.GroupMemberModel
	err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("group_id IN ? AND status = ?", groupIDs, 1).
		Find(&allMembers).Error

	if err != nil {
		l.Errorf("查询群组成员失败: %v", err)
		return nil, err
	}

	// 去重并排除自己
	seen := make(map[string]bool)
	var allMemberIDs []string

	for _, member := range allMembers {
		if member.UserID != in.UserID && !seen[member.UserID] {
			seen[member.UserID] = true
			allMemberIDs = append(allMemberIDs, member.UserID)
		}
	}

	l.Infof("获取用户群成员成功，用户ID: %s, 群成员数: %d", in.UserID, len(allMemberIDs))

	return &group_rpc.GetUserGroupMembersRes{
		MemberIDs: allMemberIDs,
	}, nil
}
