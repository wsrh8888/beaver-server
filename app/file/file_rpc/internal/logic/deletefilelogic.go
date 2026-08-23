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
	"os"

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteFileLogic) DeleteFile(in *file_rpc.DeleteFileReq) (*file_rpc.DeleteFileRes, error) {
	var file file_models.FileModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "文件不存在")
		}
		return nil, err
	}

	if err := l.svcCtx.DB.Delete(&file).Error; err != nil {
		l.Errorf("删除文件记录失败: %v", err)
		return nil, status.Error(codes.Internal, "删除文件失败")
	}

	if file.Path != "" {
		if err := os.Remove(file.Path); err != nil {
			l.Errorf("删除物理文件失败 path=%s: %v", file.Path, err)
		}
	}

	return &file_rpc.DeleteFileRes{Path: file.Path}, nil
}
