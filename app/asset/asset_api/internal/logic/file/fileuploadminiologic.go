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
	"mime/multipart"

	"beaver/app/asset/asset_api/internal/handler/common"
	"beaver/app/asset/asset_api/internal/svc"
	"beaver/app/asset/asset_api/internal/types"
	file_models "beaver/app/asset/asset_models/file"
	beaverlog "beaver/utils/beaverlog"
)

type FileUploadMinioLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 文件上传MinIO
func NewFileUploadMinioLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileUploadMinioLogic {
	return &FileUploadMinioLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("file_upload_minio", ctx),
	}
}

func (l *FileUploadMinioLogic) FileUploadMinio(req *types.FileUploadMinioReq, file multipart.File, fileHead *multipart.FileHeader, fileInfoStr string) (*types.FileUploadMinioRes, error) {
	if l.svcCtx.MinioClient == nil {
		return nil, errors.New("MinIO未配置")
	}
	// MinIO存储后端
	store := common.NewMinioStorage(l.svcCtx.MinioClient, l.svcCtx.Config.Minio.Bucket)
	base, err := common.UploadFile(l.ctx, l.svcCtx, file, fileHead, fileInfoStr, store, file_models.MinioSource)
	if err != nil {
		return nil, err
	}
	resp := &types.FileUploadMinioRes{}
	common.ConvertUploadRes(base, resp)
	return resp, nil
}
