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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetFileByMd5Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileByMd5Logic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByMd5Logic {
	return &GetFileByMd5Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileByMd5Logic) GetFileByMd5(in *file_rpc.GetFileByMd5Req) (*file_rpc.GetFileByMd5Res, error) {
	if in.Md5 == "" {
		return nil, status.Error(codes.InvalidArgument, "md5不能为空")
	}

	var file file_models.FileModel
	if err := l.svcCtx.DB.Take(&file, "md5 = ?", in.Md5).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 未命中时不返回 gRPC 错误，避免客户端拦截器把正常查重记成 fail
			return &file_rpc.GetFileByMd5Res{}, nil
		}
		l.Errorf("按 md5 查询文件失败: %v", err)
		return nil, err
	}

	return &file_rpc.GetFileByMd5Res{File: toFileItem(file)}, nil
}
