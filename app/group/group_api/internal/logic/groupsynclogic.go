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

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 群资料同步
func NewGroupSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSyncLogic {
	return &GroupSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupSyncLogic) GroupSync(req *types.GroupSyncReq) (resp *types.GroupSyncRes, err error) {
	resp = &types.GroupSyncRes{
		Groups: []types.GroupSyncDataItem{},
	}

	if len(req.Groups) == 0 {
		l.Infof("群资料同步完成，用户ID: %s, 无需同步的群组", req.UserID)
		return resp, nil
	}

	// 为每个群组查询版本大于等于本地版本的群组变更
	for _, groupReq := range req.Groups {
		var groups []group_models.GroupModel
		err = l.svcCtx.DB.Where("group_id = ? AND version >= ?", groupReq.GroupID, groupReq.Version).
			Find(&groups).Error
		if err != nil {
			l.Errorf("查询群组数据失败，群组ID: %s, 错误: %v", groupReq.GroupID, err)
			continue
		}

		for _, group := range groups {
			resp.Groups = append(resp.Groups, types.GroupSyncDataItem{
				GroupID:   group.GroupID,
				Title:     group.Title,
				Avatar:    group.Avatar,
				CreatorID: group.CreatorID,
				JoinType:  group.JoinType,
				Status:    group.Status,
				Notice:    group.Notice,
				Version:   group.Version,
				CreatedAt: time.Time(group.CreatedAt).Unix(),
				UpdatedAt: time.Time(group.UpdatedAt).Unix(),
			})
		}
	}

	l.Infof("群资料同步完成，用户ID: %s, 返回群组变化数: %d", req.UserID, len(resp.Groups))
	return resp, nil
}
