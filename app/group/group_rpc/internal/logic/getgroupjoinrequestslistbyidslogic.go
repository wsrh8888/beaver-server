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
	"fmt"
	"time"

	"beaver/app/group/group_models"
	"beaver/app/group/group_rpc/internal/svc"
	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupJoinRequestsListByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupJoinRequestsListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupJoinRequestsListByIdsLogic {
	return &GetGroupJoinRequestsListByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetGroupJoinRequestsListByIdsLogic) GetGroupJoinRequestsListByIds(in *group_rpc.GetGroupJoinRequestsListByIdsReq) (*group_rpc.GetGroupJoinRequestsListByIdsRes, error) {
	if len(in.GroupIDs) == 0 {
		l.Errorf("群组ID列表为空")
		return &group_rpc.GetGroupJoinRequestsListByIdsRes{Requests: []*group_rpc.GroupJoinRequestListById{}}, nil
	}

	// 查询指定群组ID列表中的入群申请
	var requestsData []group_models.GroupJoinRequestModel
	query := l.svcCtx.DB.Where("group_id IN (?)", in.GroupIDs)

	// 注意：Since在这里表示客户端已知的最新版本号，用于增量同步
	if in.Since > 0 {
		query = query.Where("version > ?", in.Since)
	}

	err := query.Find(&requestsData).Error
	if err != nil {
		l.Errorf("查询入群申请失败: groupIDs=%v, since=%d, error=%v", in.GroupIDs, in.Since, err)
		return nil, err
	}

	l.Infof("查询到 %d 个入群申请", len(requestsData))

	// 转换为响应格式
	var requests []*group_rpc.GroupJoinRequestListById
	for _, request := range requestsData {
		requests = append(requests, &group_rpc.GroupJoinRequestListById{
			RequestID: fmt.Sprintf("%d", request.Id), // 使用主键Id作为RequestID
			GroupID:   request.GroupID,
			UserID:    request.ApplicantUserID,
			Message:   request.Message,
			Status:    int32(request.Status),
			AppliedAt: time.Time(request.CreatedAt).UnixMilli(), // 使用CreatedAt作为申请时间
			Version:   request.Version,
		})
	}

	return &group_rpc.GetGroupJoinRequestsListByIdsRes{Requests: requests}, nil
}
