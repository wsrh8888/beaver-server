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

	filecommon "beaver/app/backend/backend_admin/internal/handler/file/common"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileListLogic {
	return &GetFileListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetFileListLogic) GetFileList(req *types.GetFileListReq) (resp *types.GetFileListRes, err error) {
	rpcRes, err := l.svcCtx.FileRpc.ListFiles(l.ctx, &file_rpc.ListFilesReq{
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		Type:     req.Type,
		Keywords: req.Keywords,
	})
	if err != nil {
		l.Errorf("获取文件列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetFileListItem, 0, len(rpcRes.List))
	for _, f := range rpcRes.List {
		list = append(list, types.GetFileListItem{
			Id:           uint(f.Id),
			FileName:     f.FileKey,
			OriginalName: f.OriginalName,
			Size:         f.Size,
			Path:         buildFileAccessURL(l.svcCtx, f),
			Md5:          f.Md5,
			Type:         f.Type,
			CreatedAt:    f.CreatedAt,
			UpdatedAt:    f.UpdatedAt,
		})
	}

	return &types.GetFileListRes{List: list, Total: rpcRes.Total}, nil
}

// buildFileAccessURL 按来源拼完整访问 URL（本地用 Domain，七牛用 Qiniu.Domain）
func buildFileAccessURL(svcCtx *svc.ServiceContext, f *file_rpc.FileItem) string {
	if f == nil {
		return ""
	}
	if f.Source == "qiniu" {
		return filecommon.BuildQiniuFileURL(svcCtx.Config.Qiniu.Domain, f.Path)
	}
	return filecommon.BuildLocalFileURL(svcCtx.Config.Domain, f.FileKey)
}
