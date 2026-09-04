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

type GroupJoinRequestListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 获取用户管理的群组申请列表
func NewGroupJoinRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupJoinRequestListLogic {
	return &GroupJoinRequestListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("group_join_request_list", ctx),
	}
}

func (l *GroupJoinRequestListLogic) GroupJoinRequestList(req *types.GroupJoinRequestListReq) (resp *types.GroupJoinRequestListRes, err error) {
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

	// 先获取用户管理的群组ID列表（作为群主或管理员）
	var managedGroupIDs []string
	err = l.svcCtx.DB.Model(&group_models.GroupMemberModel{}).
		Where("user_id = ? AND role IN (1, 2)", req.UserID).
		Pluck("group_id", &managedGroupIDs).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取用户管理的群组失败", Data: map[string]any{"userId": req.UserID, "err": err.Error()}})
		return nil, err
	}

	// 如果用户没有管理的群组，直接返回空结果
	if len(managedGroupIDs) == 0 {
		return &types.GroupJoinRequestListRes{
			List:  []types.GroupJoinRequestItem{},
			Count: 0,
		}, nil
	}

	var requests []group_models.GroupJoinRequestModel

	// 查询用户管理的所有群组的申请列表
	err = l.svcCtx.DB.Where("group_id IN (?)", managedGroupIDs).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&requests).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "查询群组申请列表失败", Data: map[string]any{"userId": req.UserID, "err": err.Error()}})
		return nil, err
	}

	// 获取总数
	var total int64
	err = l.svcCtx.DB.Model(&group_models.GroupJoinRequestModel{}).
		Where("group_id IN (?)", managedGroupIDs).
		Count(&total).Error
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取群组申请总数失败", Data: map[string]any{"userId": req.UserID, "err": err.Error()}})
		return nil, err
	}

	// 转换为响应格式
	var requestItems []types.GroupJoinRequestItem

	for _, request := range requests {
		// 这里需要查询用户信息，但由于没有用户RPC，暂时使用默认值
		// 在实际项目中，应该通过用户RPC获取用户昵称和头像
		requestItems = append(requestItems, types.GroupJoinRequestItem{
			RequestID:       request.Id,
			GroupID:         request.GroupID,
			ApplicantID:     request.ApplicantUserID,
			ApplicantName:   "用户" + request.ApplicantUserID, // 临时值，需要从用户服务获取
			ApplicantAvatar: "",                             // 临时值，需要从用户服务获取
			Message:         request.Message,
			Status:          request.Status,
			CreatedAt:       time.Time(request.CreatedAt).Unix(),
			Version:         request.Version,
		})
	}

	resp = &types.GroupJoinRequestListRes{
		List:  requestItems,
		Count: total,
	}

	l.logger.Info(model.LogMsg{Text: "获取群组申请列表完成", Data: map[string]interface{}{"userId": req.UserID, "count": len(requestItems)}})
	return resp, nil
}
