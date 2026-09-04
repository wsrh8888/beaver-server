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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type GroupMemberSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 群成员同步
func NewGroupMemberSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberSyncLogic {
	return &GroupMemberSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("group_member_sync", ctx),
	}
}

func (l *GroupMemberSyncLogic) GroupMemberSync(req *types.GroupMemberSyncReq) (resp *types.GroupMemberSyncRes, err error) {
	resp = &types.GroupMemberSyncRes{
		GroupMembers: []types.GroupMemberSyncDataItem{},
	}

	if len(req.Groups) == 0 {
		l.logger.Info(model.LogMsg{Text: "群成员同步无需群组", Data: map[string]any{"userId": req.UserID}})
		return resp, nil
	}

	// 为每个群组查询版本大于本地版本的成员变更（包括所有状态：正常、退出、被踢）
	for _, groupReq := range req.Groups {
		var members []group_models.GroupMemberModel
		// 查询 version > groupReq.Version 的成员，包括所有状态
		err = l.svcCtx.DB.Where("group_id = ? AND version >= ?", groupReq.GroupID, groupReq.Version).
			Find(&members).Error
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "查询群成员数据失败", Data: map[string]any{"groupId": groupReq.GroupID, "err": err.Error()}})
			continue
		}

		for _, member := range members {
			resp.GroupMembers = append(resp.GroupMembers, types.GroupMemberSyncDataItem{
				GroupID:   member.GroupID,
				UserID:    member.UserID,
				Role:      member.Role,
				Status:    member.Status,
				JoinTime:  member.JoinTime.Unix(),
				Version:   member.Version,
				CreatedAt: time.Time(member.CreatedAt).Unix(),
				UpdatedAt: time.Time(member.UpdatedAt).Unix(),
			})
		}
	}

	l.logger.Info(model.LogMsg{Text: "群成员同步完成", Data: map[string]interface{}{"userId": req.UserID, "count": len(resp.GroupMembers)}})
	return resp, nil
}
