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
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupsListByIdsLogic {
	return &GetGroupsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupsListByIdsLogic) GetGroupsListByIds(in *group_rpc.GetGroupsListByIdsReq) (*group_rpc.GetGroupsListByIdsRes, error) {
	if len(in.GroupIDs) == 0 {
		l.Errorf("群组ID列表为空")
		return &group_rpc.GetGroupsListByIdsRes{Groups: []*group_rpc.GroupListById{}}, nil
	}
	// 查询指定群组ID列表中的群组资料 - 只查询需要的字段
	var groupsData []group_models.GroupModel
	query := l.svcCtx.DB.Select("group_id, title, avatar, notice, version, created_at, updated_at").
		Where("group_id IN (?)", in.GroupIDs)

	// 注意：Since在这里表示客户端已知的最新更新时间（时间戳），用于增量同步
	// 如果Since > 0，只返回更新时间大于Since的群组（有变更的群组）
	if in.Since > 0 {
		// Since是毫秒时间戳，需要转换为time.Time进行比较
		sinceTime := time.UnixMilli(in.Since)
		query = query.Where("updated_at > ?", sinceTime)
	}

	err := query.Find(&groupsData).Error
	if err != nil {
		l.Errorf("查询群组资料失败: groupIDs=%v, since=%d, error=%v", in.GroupIDs, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个群组资料", len(groupsData))

	// 转换为响应格式
	var groups []*group_rpc.GroupListById
	for _, group := range groupsData {
		groups = append(groups, &group_rpc.GroupListById{
			GroupID:     group.GroupID,
			Name:        group.Title,
			Avatar:      group.Avatar,
			Description: group.Notice,
			Version:     group.Version,
			CreatedAt:   time.Time(group.CreatedAt).UnixMilli(),
			UpdatedAt:   time.Time(group.UpdatedAt).UnixMilli(),
		})
	}

	return &group_rpc.GetGroupsListByIdsRes{Groups: groups}, nil
}
