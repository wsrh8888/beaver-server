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
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetFileByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewGetFileByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileByIdLogic {
	return &GetFileByIdLogic{ctx: ctx, svcCtx: svcCtx, logger: beaverlog.New("get_file_by_id", ctx)}
}

func (l *GetFileByIdLogic) GetFileById(in *file_rpc.GetFileByIdReq) (*file_rpc.GetFileByIdRes, error) {
	var file file_models.FileModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&file).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "文件不存在")
		}
		l.logger.Error(model.LogMsg{
			Text: "按ID查询文件失败",
			Data: map[string]any{"id": in.Id, "err": err.Error()},
		})
		return nil, err
	}
	return &file_rpc.GetFileByIdRes{File: toFileItem(file)}, nil
}
