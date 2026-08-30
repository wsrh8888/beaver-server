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
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileDetailLogic {
	return &GetFileDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetFileDetailLogic) GetFileDetail(req *types.GetFileDetailReq) (resp *types.GetFileDetailRes, err error) {
	rpcRes, err := l.svcCtx.FileRpc.GetFileById(l.ctx, &file_rpc.GetFileByIdReq{Id: uint64(req.Id)})
	if err != nil {
		l.Errorf("获取文件详情失败: %v", err)
		return nil, err
	}

	f := rpcRes.File
	return &types.GetFileDetailRes{
		Id:           uint(f.Id),
		FileName:     f.FileKey,
		OriginalName: f.OriginalName,
		Size:         f.Size,
		Path:         buildFileAccessURL(l.svcCtx, f),
		Md5:          f.Md5,
		Type:         f.Type,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}, nil
}
