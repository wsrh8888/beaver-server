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

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileDetailLogic {
	return &GetFileDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通过fileName查询文件详情
func (l *GetFileDetailLogic) GetFileDetail(in *file_rpc.GetFileDetailReq) (*file_rpc.GetFileDetailRes, error) {
	var file file_models.FileModel

	// 通过fileName查询文件信息
	err := l.svcCtx.DB.Take(&file, "file_key = ?", in.FileKey).Error
	if err != nil {
		logx.Errorf("查询文件失败: %s", err.Error())
		return nil, errors.New("文件不存在")
	}

	// 返回文件详情
	return &file_rpc.GetFileDetailRes{
		FileKey:      file.FileKey,
		OriginalName: file.OriginalName,
		Size:         file.Size,
		Path:         file.Path,
		Md5:          file.Md5,
		Type:         file.Type,
		CreatedAt:    file.CreatedAt.String(),
		UpdatedAt:    file.UpdatedAt.String(),
	}, nil
}
