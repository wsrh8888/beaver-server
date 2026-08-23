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
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMomentCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentCommentListLogic {
	return &GetMomentCommentListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMomentCommentListLogic) GetMomentCommentList(req *types.GetMomentCommentListReq) (resp *types.GetMomentCommentListRes, err error) {
	if req.MomentId == "" {
		return nil, errors.New("动态ID不能为空")
	}

	rpcRes, err := l.svcCtx.MomentRpc.ListMomentComments(l.ctx, &moment_rpc.ListMomentCommentsReq{
		MomentId: req.MomentId,
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
	})
	if err != nil {
		l.Errorf("获取动态评论列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetMomentCommentListItem, 0, len(rpcRes.List))
	for _, c := range rpcRes.List {
		list = append(list, types.GetMomentCommentListItem{
			CommentId: c.CommentId,
			MomentId:  c.MomentId,
			UserId:    c.UserId,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}
	return &types.GetMomentCommentListRes{List: list, Total: rpcRes.Total}, nil
}
