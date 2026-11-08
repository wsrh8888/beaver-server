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

	knowledge_models "beaver/app/platform/platform_models/knowledge"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListDocumentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListDocumentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListDocumentsLogic {
	return &ListDocumentsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListDocumentsLogic) ListDocuments(in *platform_rpc.ListDocumentsReq) (*platform_rpc.ListDocumentsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	size := int(in.PageSize)
	if size <= 0 || size > 200 {
		size = 50
	}
	// 软删哨兵过滤：只列未删文档，与 worker 捞取口径一致。
	tx := l.svcCtx.DB.Model(&knowledge_models.KnowledgeDocument{}).Where("deleted_at = ?", zeroDeletedAt)
	if trim(in.CollectionId) != "" {
		tx = tx.Where("collection_id = ?", in.CollectionId)
	}
	if trim(in.IndexStatus) != "" {
		tx = tx.Where("index_status = ?", in.IndexStatus)
	}
	if in.Status != 0 {
		tx = tx.Where("status = ?", in.Status)
	}
	var total int64
	tx.Count(&total)

	var rows []knowledge_models.KnowledgeDocument
	if err := tx.Order("id DESC").Limit(size).Offset((page - 1) * size).Find(&rows).Error; err != nil {
		l.Errorf("列出知识文档失败: %v", err)
		return nil, err
	}
	items := make([]*platform_rpc.KnowledgeDocumentItem, 0, len(rows))
	for i := range rows {
		items = append(items, toDocumentItem(&rows[i]))
	}
	return &platform_rpc.ListDocumentsRes{Total: total, Documents: items}, nil
}
