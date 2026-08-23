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

type GetUserGroupIDsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserGroupIDsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserGroupIDsLogic {
	return &GetUserGroupIDsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserGroupIDsLogic) GetUserGroupIDs(in *group_rpc.GetUserGroupIDsReq) (*group_rpc.GetUserGroupIDsRes, error) {
	// 获取用户加入的所有群组ID
	var members []group_models.GroupMemberModel

	err := l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ? AND status = ?", in.UserID, 1).
		Find(&members).Error

	if err != nil {
		l.Errorf("查询用户群组ID失败: %v", err)
		return nil, err
	}

	// 提取群组ID列表
	groupIDs := make([]string, 0, len(members))
	for _, member := range members {
		groupIDs = append(groupIDs, member.GroupID)
	}

	l.Infof("获取用户群组ID成功，用户ID: %s, 群组数: %d", in.UserID, len(groupIDs))

	return &group_rpc.GetUserGroupIDsRes{
		GroupIDs: groupIDs,
	}, nil
}
