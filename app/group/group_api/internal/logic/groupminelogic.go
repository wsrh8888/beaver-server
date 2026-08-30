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

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取我加入的群组列表
func NewGroupMineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMineLogic {
	return &GroupMineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMineLogic) GroupMine(req *types.GroupMineReq) (resp *types.GroupMineRes, err error) {
	var groupMembers []group_models.GroupMemberModel

	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 查询用户加入的群组
	err = l.svcCtx.DB.Where("user_id = ? AND status = ?", req.UserID, 1).
		Order("join_time DESC").
		Offset(offset).
		Limit(limit).
		Find(&groupMembers).Error
	if err != nil {
		l.Errorf("查询用户群组失败: %v", err)
		return nil, err
	}

	// 获取总数
	var total int64
	err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ? AND status = ?", req.UserID, 1).
		Count(&total).Error
	if err != nil {
		l.Errorf("获取用户群组总数失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var groupItems []types.GroupMineItem

	for _, member := range groupMembers {
		// 查询群组信息
		var group group_models.GroupModel
		err = l.svcCtx.DB.Where("group_id = ?", member.GroupID).First(&group).Error
		if err != nil {
			l.Errorf("查询群组信息失败，群组ID: %s, 错误: %v", member.GroupID, err)
			continue
		}

		// 获取群成员数量
		var memberCount int64
		err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
			Where("group_id = ? AND status = ?", group.GroupID, 1).
			Count(&memberCount).Error
		if err != nil {
			l.Errorf("获取群成员数量失败，群组ID: %s, 错误: %v", group.GroupID, err)
			continue
		}

		groupItems = append(groupItems, types.GroupMineItem{
			GroupID:        group.GroupID,
			Title:          group.Title,
			Avatar:         group.Avatar,
			MemberCount:    int(memberCount),
			ConversationID: group.GroupID, // 群组ID作为会话ID
			Version:        group.Version,
		})
	}

	resp = &types.GroupMineRes{
		List:  groupItems,
		Count: int(total),
	}

	l.Infof("获取用户群组列表完成，用户ID: %s, 返回群组数: %d", req.UserID, len(groupItems))
	return resp, nil
}
