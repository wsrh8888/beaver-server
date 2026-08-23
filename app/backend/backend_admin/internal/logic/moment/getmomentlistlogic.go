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

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_rpc/types/moment_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMomentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentListLogic {
	return &GetMomentListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMomentListLogic) GetMomentList(req *types.GetMomentListReq) (resp *types.GetMomentListRes, err error) {
	rpcRes, err := l.svcCtx.MomentRpc.ListMoments(l.ctx, &moment_rpc.ListMomentsReq{
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		UserId:   req.UserId,
		Keywords: req.Keywords,
	})
	if err != nil {
		l.Errorf("获取动态列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetMomentListItem, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		files := make([]types.GetMomentListFileInfo, 0, len(m.Files))
		for _, f := range m.Files {
			files = append(files, types.GetMomentListFileInfo{FileName: f.FileKey})
		}
		list = append(list, types.GetMomentListItem{
			MomentId:  m.MomentId,
			UserId:    m.UserId,
			Content:   m.Content,
			Files:     files,
			IsDeleted: m.IsDeleted,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		})
	}
	return &types.GetMomentListRes{List: list, Total: rpcRes.Total}, nil
}
