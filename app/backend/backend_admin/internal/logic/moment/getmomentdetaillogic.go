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

type GetMomentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMomentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentDetailLogic {
	return &GetMomentDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetMomentDetailLogic) GetMomentDetail(req *types.GetMomentDetailReq) (resp *types.GetMomentDetailRes, err error) {
	if req.MomentId == "" {
		return nil, errors.New("动态ID不能为空")
	}

	rpcRes, err := l.svcCtx.MomentRpc.ListMoments(l.ctx, &moment_rpc.ListMomentsReq{
		MomentId: req.MomentId,
		Page:     1,
		PageSize: 1,
	})
	if err != nil {
		l.Errorf("获取动态详情失败: %v", err)
		return nil, err
	}
	if len(rpcRes.List) == 0 {
		return nil, errors.New("动态不存在")
	}

	m := rpcRes.List[0]
	files := make([]types.GetMomentDetailFileInfo, 0, len(m.Files))
	for _, f := range m.Files {
		files = append(files, types.GetMomentDetailFileInfo{FileName: f.FileKey})
	}
	return &types.GetMomentDetailRes{
		MomentId:  m.MomentId,
		UserId:    m.UserId,
		Content:   m.Content,
		Files:     files,
		IsDeleted: m.IsDeleted,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
