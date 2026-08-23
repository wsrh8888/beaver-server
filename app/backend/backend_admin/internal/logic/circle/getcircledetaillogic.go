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
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

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
	if req.CircleId == "" {
		return nil, errors.New("圈子ID不能为空")
	}

	rpcRes, err := l.svcCtx.CircleRpc.GetCircleList(l.ctx, &circle_rpc.GetCircleListReq{
		CircleId: req.CircleId,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		l.Errorf("获取圈子详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("圈子不存在")
	}

	c := mapCircleItem(rpcRes.List[0])
	return &types.GetCircleDetailRes{
		CircleId:    c.CircleId,
		Name:        c.Name,
		Description: c.Description,
		Avatar:      c.Avatar,
		CreatorId:   c.CreatorId,
		JoinType:    c.JoinType,
		MemberCount: c.MemberCount,
		PostCount:   c.PostCount,
		IsDeleted:   c.IsDeleted,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}
