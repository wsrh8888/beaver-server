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

package circle

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetJoinRequestListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetJoinRequestListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetJoinRequestListLogic {
	return &GetJoinRequestListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetJoinRequestListLogic) GetJoinRequestList(req *types.GetJoinRequestListReq) (resp *types.GetJoinRequestListRes, err error) {
	// 权限校验
	var operator circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&operator).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if operator.Role > 2 {
		return nil, fmt.Errorf("仅圈主和管理员可查看申请列表")
	}

	var total int64
	var list []circle_models.CircleJoinRequestModel
	l.svcCtx.DB.Model(&circle_models.CircleJoinRequestModel{}).
		Where("circle_id = ? AND status = 0", req.CircleID).
		Count(&total)
	l.svcCtx.DB.Where("circle_id = ? AND status = 0", req.CircleID).
		Order("created_at DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&list)

	if len(list) == 0 {
		return &types.GetJoinRequestListRes{Count: total, List: []types.GetJoinRequestListItem{}}, nil
	}

	// 批量拉取用户信息
	userIDs := make([]string, 0, len(list))
	for _, r := range list {
		userIDs = append(userIDs, r.UserID)
	}
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	items := make([]types.GetJoinRequestListItem, 0, len(list))
	for _, r := range list {
		item := types.GetJoinRequestListItem{
			RequestID: fmt.Sprintf("%d", r.Id),
			UserID:    r.UserID,
			Reason:    r.Reason,
			CreatedAt: r.CreatedAt.String(),
		}
		if userResp != nil {
			if info := userResp.UserInfo[r.UserID]; info != nil {
				item.UserName = info.NickName
				item.Avatar = info.Avatar
			}
		}
		items = append(items, item)
	}

	return &types.GetJoinRequestListRes{Count: total, List: items}, nil
}
