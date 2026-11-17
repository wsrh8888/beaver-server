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

	"beaver/app/asset/asset_api/internal/svc"
	"beaver/app/asset/asset_api/internal/types"
	assettypes "beaver/app/asset/asset_rpc/pb"
	beaverlog "beaver/utils/beaverlog"
)

type FileByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 队间：按记录ID取文件
func NewFileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileByIdLogic {
	return &FileByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("FileById", ctx),
	}
}

// FileById 只做协议转换，业务规则全部留在 asset_rpc。
func (l *FileByIdLogic) FileById(req *types.GetFileByIdReq) (resp *types.GetFileByIdRes, err error) {
	rpcReq := &assettypes.GetFileByIdReq{Id: req.ID}
	rpcRes, err := l.svcCtx.FileRpc.GetFileById(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.GetFileByIdRes{File: &types.FileItem{ID: rpcRes.File.Id, FileKey: rpcRes.File.FileKey, OriginalName: rpcRes.File.OriginalName, Size: rpcRes.File.Size, Path: rpcRes.File.Path, Md5: rpcRes.File.Md5, Type: rpcRes.File.Type, Source: rpcRes.File.Source, CreatedAt: rpcRes.File.CreatedAt, UpdatedAt: rpcRes.File.UpdatedAt}}, nil
}
