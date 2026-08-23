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
	"strings"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleDetailLogic {
	return &GetCircleDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleDetailLogic) GetCircleDetail(req *types.GetCircleDetailReq) (resp *types.GetCircleDetailRes, err error) {
	var circle circle_models.CircleModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", req.CircleID).First(&circle).Error; err != nil {
		return nil, fmt.Errorf("圈子不存在")
	}

	role := int8(0)
	var member circle_models.CircleMemberModel
	if l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error == nil {
		role = member.Role
	}

	resp = &types.GetCircleDetailRes{
		CircleID:    circle.CircleID,
		Name:        circle.Name,
		Description: circle.Description,
		Avatar:      circle.Avatar,
		JoinType:    circle.JoinType,
		CreatorID:   circle.CreatorID,
		MemberCount: countMembers(l.svcCtx.DB, req.CircleID),
		PostCount:   countPosts(l.svcCtx.DB, req.CircleID),
		Role:        role,
		CreatedAt:   circle.CreatedAt.String(),
	}
	if role > 0 {
		var invite circle_models.CircleInviteModel
		if e := l.svcCtx.DB.Where("circle_id = ? AND status = 1", circle.CircleID).Order("id asc").First(&invite).Error; e == nil {
			domain := strings.TrimRight(strings.TrimSpace(l.svcCtx.Config.Domain), "/")
			if domain != "" && invite.Token != "" {
				resp.InviteUrl = fmt.Sprintf("%s/api/circle/v1/circle/invite_code?code=%s", domain, invite.Token)
			}
		}
	}
	return resp, nil
}
