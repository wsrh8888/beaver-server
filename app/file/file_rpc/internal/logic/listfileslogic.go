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

	"beaver/app/file/file_models"
	"beaver/app/file/file_rpc/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListFilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFilesLogic {
	return &ListFilesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListFilesLogic) ListFiles(in *file_rpc.ListFilesReq) (*file_rpc.ListFilesRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&file_models.FileModel{})
	if in.Type != "" {
		db = db.Where("type = ?", in.Type)
	}
	if in.Keywords != "" {
		db = db.Where("file_key LIKE ? OR original_name LIKE ?", "%"+in.Keywords+"%", "%"+in.Keywords+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计文件失败: %v", err)
		return nil, err
	}

	var list []file_models.FileModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		l.Errorf("查询文件列表失败: %v", err)
		return nil, err
	}

	items := make([]*file_rpc.FileItem, 0, len(list))
	for _, f := range list {
		items = append(items, toFileItem(f))
	}

	return &file_rpc.ListFilesRes{Total: total, List: items}, nil
}
