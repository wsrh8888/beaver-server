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

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCircleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleListLogic {
	return &GetCircleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCircleListLogic) GetCircleList(in *circle_rpc.GetCircleListReq) (*circle_rpc.GetCircleListRes, error) {
	query := l.svcCtx.DB.Model(&circle_models.CircleModel{}).Where("is_deleted = false")

	if in.CircleId != "" {
		query = query.Where("circle_id = ?", in.CircleId)
	}
	if in.UserId != "" {
		var memberCircleIDs []string
		l.svcCtx.DB.Model(&circle_models.CircleMemberModel{}).
			Where("user_id = ?", in.UserId).
			Pluck("circle_id", &memberCircleIDs)
		query = query.Where("circle_id IN ?", memberCircleIDs)
	}
	if in.Keywords != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+in.Keywords+"%", "%"+in.Keywords+"%")
	}

	var total int64
	query.Count(&total)

	var circles []circle_models.CircleModel
	query.Order("created_at DESC").
		Offset(int((in.Page - 1) * in.PageSize)).
		Limit(int(in.PageSize)).
		Find(&circles)

	circleIDs := make([]string, 0, len(circles))
	for _, c := range circles {
		circleIDs = append(circleIDs, c.CircleID)
	}
	memberCountMap := countMembersByCircleIDs(l.svcCtx.DB, circleIDs)
	postCountMap := countPostsByCircleIDs(l.svcCtx.DB, circleIDs)

	list := make([]*circle_rpc.CircleItem, 0, len(circles))
	for _, c := range circles {
		list = append(list, &circle_rpc.CircleItem{
			CircleId:    c.CircleID,
			Name:        c.Name,
			Description: c.Description,
			Avatar:      c.Avatar,
			CreatorId:   c.CreatorID,
			JoinType:    int32(c.JoinType),
			MemberCount: memberCountMap[c.CircleID],
			PostCount:   postCountMap[c.CircleID],
			IsDeleted:   c.IsDeleted,
			CreatedAt:   c.CreatedAt.String(),
			UpdatedAt:   c.UpdatedAt.String(),
		})
	}

	return &circle_rpc.GetCircleListRes{Total: total, List: list}, nil
}
