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

	knowledge_models "beaver/app/knowledge/knowledge_models"
	"beaver/app/knowledge/knowledge_rpc/internal/svc"
	"beaver/app/knowledge/knowledge_rpc/types/knowledge_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCollectionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCollectionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCollectionsLogic {
	return &ListCollectionsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListCollectionsLogic) ListCollections(in *knowledge_rpc.ListCollectionsReq) (*knowledge_rpc.ListCollectionsRes, error) {
	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	size := int(in.PageSize)
	if size <= 0 || size > 200 {
		size = 50
	}
	tx := l.svcCtx.DB.Model(&knowledge_models.KnowledgeCollection{})
	if in.Status != 0 {
		tx = tx.Where("status = ?", in.Status)
	}
	var total int64
	tx.Count(&total)

	var rows []knowledge_models.KnowledgeCollection
	if err := tx.Order("sort ASC, id ASC").Limit(size).Offset((page - 1) * size).Find(&rows).Error; err != nil {
		l.Errorf("列出知识集合失败: %v", err)
		return nil, err
	}
	items := make([]*knowledge_rpc.KnowledgeCollectionItem, 0, len(rows))
	for i := range rows {
		items = append(items, toCollectionItem(&rows[i]))
	}
	return &knowledge_rpc.ListCollectionsRes{Total: total, Collections: items}, nil
}
