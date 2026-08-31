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

type GroupSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 群资料同步
func NewGroupSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSyncLogic {
	return &GroupSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("group_sync", ctx),
	}
}

func (l *GroupSyncLogic) GroupSync(req *types.GroupSyncReq) (resp *types.GroupSyncRes, err error) {
	resp = &types.GroupSyncRes{
		Groups: []types.GroupSyncDataItem{},
	}

	if len(req.Groups) == 0 {
		l.logger.Info(model.LogMsg{Text: "群资料同步无需群组", Data: map[string]any{"userId": req.UserID}})
		return resp, nil
	}

	// 为每个群组查询版本大于等于本地版本的群组变更
	for _, groupReq := range req.Groups {
		var groups []group_models.GroupModel
		err = l.svcCtx.DB.Where("group_id = ? AND version >= ?", groupReq.GroupID, groupReq.Version).
			Find(&groups).Error
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "查询群组数据失败", Data: map[string]any{"groupId": groupReq.GroupID, "err": err.Error()}})
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

	l.logger.Info(model.LogMsg{Text: "群资料同步完成", Data: map[string]interface{}{"userId": req.UserID, "count": len(resp.Groups)}})
	return resp, nil
}
