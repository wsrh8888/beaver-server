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

type CircleSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCircleSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CircleSyncLogic {
	return &CircleSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CircleSyncLogic) CircleSync(req *types.CircleSyncReq) (resp *types.CircleSyncRes, err error) {
	var members []circle_models.CircleMemberModel
	l.svcCtx.DB.Where("user_id = ?", req.UserID).Find(&members)

	if len(members) == 0 {
		return &types.CircleSyncRes{List: []types.CircleSyncItem{}}, nil
	}

	circleIDs := make([]string, 0, len(members))
	roleMap := make(map[string]int8)
	for _, m := range members {
		circleIDs = append(circleIDs, m.CircleID)
		roleMap[m.CircleID] = m.Role
	}

	var circles []circle_models.CircleModel
	l.svcCtx.DB.Where("circle_id IN ? AND version > ?", circleIDs, req.Version).Find(&circles)

	syncIDs := make([]string, 0, len(circles))
	for _, c := range circles {
		syncIDs = append(syncIDs, c.CircleID)
	}
	memberCountMap := countMembersByCircleIDs(l.svcCtx.DB, syncIDs)

	items := make([]types.CircleSyncItem, 0, len(circles))
	for _, c := range circles {
		role := roleMap[c.CircleID]
		if c.IsDeleted {
			role = 0
		}
		items = append(items, types.CircleSyncItem{
			CircleID:    c.CircleID,
			Name:        c.Name,
			Avatar:      c.Avatar,
			MemberCount: memberCountMap[c.CircleID],
			Role:        role,
			Version:     c.Version,
		})
	}

	return &types.CircleSyncRes{List: items}, nil
}
