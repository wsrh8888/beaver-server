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
	"mime/multipart"

	"beaver/app/file/file_api/internal/handler/common"
	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
)

type FileUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 文件通用上传(按yaml StorageType配置路由到local/qiniu/minio)
func NewFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadLogic {
	return &FileUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("file_upload", ctx),
	}
}

func (l *FileUploadLogic) FileUpload(req *types.FileReq, file multipart.File, fileHead *multipart.FileHeader, fileInfoStr string) (*types.FileRes, error) {
	// 按配置选择默认存储后端
	store, source, err := common.DefaultStorage(l.svcCtx)
	if err != nil {
		return nil, err
	}
	return common.UploadFile(l.ctx, l.svcCtx, file, fileHead, fileInfoStr, store, source)
}
