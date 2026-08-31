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

type GroupJoinRequestSyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 入群申请同步
func NewGroupJoinRequestSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinRequestSyncLogic {
	return &GroupJoinRequestSyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("group_join_request_sync", ctx),
	}
}

func (l *GroupJoinRequestSyncLogic) GroupJoinRequestSync(req *types.GroupJoinRequestSyncReq) (resp *types.GroupJoinRequestSyncRes, err error) {
	resp = &types.GroupJoinRequestSyncRes{
		GroupJoinRequests: []types.GroupJoinRequestSyncDataItem{},
	}

	if len(req.Groups) == 0 {
		l.logger.Info(model.LogMsg{Text: "入群申请同步无需群组", Data: map[string]any{"userId": req.UserID}})
		return resp, nil
	}

	// 为每个群组查询版本变化的数据
	for _, groupReq := range req.Groups {
		var requests []group_models.GroupJoinRequestModel
		err = l.svcCtx.DB.Where("group_id = ? AND version >= ?", groupReq.GroupID, groupReq.Version).
			Find(&requests).Error
		if err != nil {
			l.logger.Error(model.LogMsg{Text: "查询入群申请数据失败", Data: map[string]any{"groupId": groupReq.GroupID, "err": err.Error()}})
			continue
		}

		for _, request := range requests {
			resp.GroupJoinRequests = append(resp.GroupJoinRequests, types.GroupJoinRequestSyncDataItem{
				GroupID:         request.GroupID,
				ApplicantUserID: request.ApplicantUserID,
				Message:         request.Message,
				Status:          request.Status,
				HandledBy:       request.HandledBy,
				HandledAt: func() int64 {
					if request.HandledAt != nil {
						return request.HandledAt.Unix()
					}
					return 0
				}(),
				Version:   request.Version,
				CreatedAt: time.Time(request.CreatedAt).Unix(),
				UpdatedAt: time.Time(request.UpdatedAt).Unix(),
			})
		}
	}

	l.logger.Info(model.LogMsg{Text: "入群申请同步完成", Data: map[string]interface{}{"userId": req.UserID, "count": len(resp.GroupJoinRequests)}})
	return resp, nil
}
