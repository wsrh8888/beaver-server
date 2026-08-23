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

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchCircleLogic {
	return &SearchCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchCircleLogic) SearchCircle(req *types.SearchCircleReq) (resp *types.SearchCircleRes, err error) {
	var total int64
	var circles []circle_models.CircleModel

	query := l.svcCtx.DB.Model(&circle_models.CircleModel{}).Where("is_deleted = false")
	if req.Keywords != "" {
		query = query.Where("name LIKE ? OR description LIKE ?",
			"%"+req.Keywords+"%", "%"+req.Keywords+"%")
	}
	query.Count(&total)
	query.Order("created_at DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&circles)

	if len(circles) == 0 {
		return &types.SearchCircleRes{Count: total, List: []types.SearchCircleItem{}}, nil
	}

	circleIDs := make([]string, 0, len(circles))
	for _, c := range circles {
		circleIDs = append(circleIDs, c.CircleID)
	}
	var members []circle_models.CircleMemberModel
	l.svcCtx.DB.Where("circle_id IN ? AND user_id = ?", circleIDs, req.UserID).Find(&members)
	roleMap := make(map[string]int8)
	for _, m := range members {
		roleMap[m.CircleID] = m.Role
	}
	memberCountMap := countMembersByCircleIDs(l.svcCtx.DB, circleIDs)

	items := make([]types.SearchCircleItem, 0, len(circles))
	for _, c := range circles {
		items = append(items, types.SearchCircleItem{
			CircleID:    c.CircleID,
			Name:        c.Name,
			Description: c.Description,
			Avatar:      c.Avatar,
			MemberCount: memberCountMap[c.CircleID],
			JoinType:    c.JoinType,
			Role:        roleMap[c.CircleID],
		})
	}

	return &types.SearchCircleRes{Count: total, List: items}, nil
}
