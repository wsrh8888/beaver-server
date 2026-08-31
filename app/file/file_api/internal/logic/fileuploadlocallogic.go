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

	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type FileUploadLocalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 文件上传本地
func NewFileUploadLocalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLocalLogic {
	return &FileUploadLocalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("file_upload_local", ctx),
	}
}

func (l *FileUploadLocalLogic) FileUploadLocal(req *types.FileReq) (resp *types.FileRes, err error) {
	// todo: add your logic here and delete this line

	return &types.FileRes{}, nil
}
